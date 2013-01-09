package aoi

import (
	"fmt"
	"math"
)

type Rectangle struct {
	Width  int
	Height int
}

type Position struct {
	X int
	Y int
}

type MapConfig struct {
	*Rectangle
	ID int
}

type TowerConfig struct {
	*Rectangle
}

type Config struct {
	M     *MapConfig
	T     *TowerConfig
	Limit int
}

type TowerAOI struct {
	*Config
	width      int
	height     int
	towers     [][]*Tower
	rangeLimit int
	max        *Position
}

func (t *TowerAOI) Init(config *Config) {
	t.Config = config

	if t.Config.Limit != 0 {
		t.rangeLimit = t.Config.Limit
	} else {
		t.rangeLimit = 5
	}

	t.width = config.T.Width
	t.height = config.T.Height

	var m int = t.Config.M.Width / t.Config.T.Width
	var n int = t.Config.M.Height / t.Config.T.Height

	t.max = &Position{
		X: m - 1,
		Y: n - 1,
	}

	t.towers = make([][]*Tower, m)
	for i := 0; i < m; i++ {
		t.towers[i] = make([]*Tower, n)
		for j := 0; j < n; j++ {
			t.towers[i][j] = CreateTower()
		}
	}
}

// Get given type object ids from tower aoi by range and types
// @param pos {Object} The pos to find objects
// @param range {Number} The range to find the object, in tower aoi, it means the tower number from the pos 
// @param types {Array} The types of the object need to find
func (t *TowerAOI) GetIdsByRange(pos *Position, r int, types []int) map[int]map[int64]int64 {
	if !t.checkPos(pos) || r < 0 || r > t.rangeLimit {
		return nil
	}

	p := t.transPos(pos)
	if p == nil {
		fmt.Println("p value is nil.")
	}
	start, end := t.getPosLimit(p, r, t.max)

	result := make(map[int]map[int64]int64)
	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			result = t.addMapByTypes(result, t.towers[i][j].GetIdsByTypes(types), types)
		}
	}
	return result
}

// Get all object ids from tower aoi by pos and range
func (t *TowerAOI) GetIdsByPos(pos *Position, r int) []int64 {
	if t.checkPos(pos) || r < 0 {
		return nil
	}

	if r > 5 {
		r = 5
	}

	p := t.transPos(pos)
	start, end := t.getPosLimit(p, r, t.max)
	var result []int64
	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			result = addMap(result, t.towers[i][j].GetIds())
		}
	}
	return result
}

// Add an object to tower aoi at given pos
func (t *TowerAOI) AddObject(obj *Object, pos *Position) bool {
	if t.checkPos(pos) {
		p := t.transPos(pos)
		t.towers[p.X][p.Y].Add(obj)
		// emit('add', {id: obj.id, type:obj.type, watchers:this.towers[p.x][p.y].watchers});
		return true
	}
	return false
}

// Remove object from aoi module
func (t *TowerAOI) RemoveObject(obj *Object, pos *Position) bool {
	if t.checkPos(pos) {
		p := t.transPos(pos)
		t.towers[p.X][p.Y].Remove(obj)
		// this.emit('remove', {id: obj.id, type:obj.type, watchers:this.towers[p.x][p.y].watchers});
		return true
	}
	return false
}

func (t *TowerAOI) UpdateObject(obj *Object, oldPos *Position, newPos *Position) bool {
	if !t.checkPos(oldPos) || !t.checkPos(newPos) {
		return false
	}

	p1 := t.transPos(oldPos)
	p2 := t.transPos(newPos)

	if p1.X == p2.X && p1.Y == p2.Y {
		return true
	} else {
		if t.towers[p1.X] == nil || t.towers[p2.X] == nil {
			fmt.Printf("AOI pos error ! oldPos : %v, newPos : %v, p1 : %v, p2 : %v\n", oldPos, newPos, p1, p2)
			return false
		}

		oldTower := t.towers[p1.X][p1.Y]
		newTower := t.towers[p2.X][p2.Y]

		oldTower.Remove(obj)
		newTower.Add(obj)

		// this.emit('update', {id: obj.id, type:obj.type, oldWatchers:oldTower.watchers, newWatchers:newTower.watchers})
		return true
	}
	return false
}

// Check if the pos is valid;
// @return {Boolean} Test result
func (t *TowerAOI) checkPos(pos *Position) bool {
	if pos == nil {
		return false
	}
	if pos.X < 0 || pos.Y < 0 || pos.X >= t.M.Width || pos.Y >= t.M.Height {
		return false
	}
	return true
}

// Trans the absolut pos to tower pos. For example : (210, 110} -> (1, 0), for tower width 200, height 200
func (t *TowerAOI) transPos(pos *Position) *Position {
	return &Position{
		X: int(math.Floor(float64(pos.X) / float64(t.Config.T.Width))),
		Y: int(math.Floor(float64(pos.Y) / float64(t.Config.T.Height))),
	}
}

func (t *TowerAOI) AddWatcher(watcher *Object, pos *Position, r int) {
	if r < 0 {
		return
	}

	if r > 5 {
		r = 5
	}

	p := t.transPos(pos)
	start, end := t.getPosLimit(p, r, t.max)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			t.towers[i][j].AddWatcher(watcher)
		}
	}
}

func (t *TowerAOI) RemoveWatcher(watcher *Object, pos *Position, r int) {
	if r < 0 {
		return
	}

	if r > 5 {
		r = 5
	}

	p := t.transPos(pos)
	start, end := t.getPosLimit(p, r, t.max)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			t.towers[i][j].RemoveWatcher(watcher)
		}
	}
}

func (t *TowerAOI) UpdateWatcher(watcher *Object, oldPos *Position, newPos *Position, oldRange int, newRange int) bool {
	if !t.checkPos(oldPos) || !t.checkPos(newPos) {
		return false
	}

	p1 := t.transPos(oldPos)
	p2 := t.transPos(newPos)

	if p1.X == p2.X && p1.Y == p2.Y {
		return true
	} else {
		if oldRange < 0 || newRange < 0 {
			return false
		}

		if oldRange > 5 {
			oldRange = 5
		}
		if newRange > 5 {
			newRange = 5
		}

		removeTowers, addTowers, unChangeTowers := t.getChangedTowers(p1, p2, oldRange, newRange, t.towers, t.max)
		addObjs := make([]int64, 1)
		removeObjs := make([]int64, 1)

		for i := 0; i < len(addTowers); i++ {
			addTowers[i].AddWatcher(watcher)
			ids := addTowers[i].GetIds()
			addObjs = addMap(addObjs, ids)
		}

		for i := 0; i < len(removeTowers); i++ {
			removeTowers[i].RemoveWatcher(watcher)
			ids := removeTowers[i].GetIds()
			removeObjs = addMap(removeObjs, ids)
		}

		fmt.Println("unChangeTowers: %v ?????", unChangeTowers)

		//this.emit('updateWatcher', {id: watcher.id, type:watcher.type, addObjs: addObjs, removeObjs:removeObjs});
		return true
	}
	return false
}

// Get changed towers for girven pos
// @param p1 {Object} The origin position
// @param p2 {Object} The now position
// @param oldRange {Number} The old range
// @param newRange {Number} The new range
// @param towers {Object} All towers of the aoi
// @param max {Object} The position limit of the towers
func (t *TowerAOI) getChangedTowers(p1 *Position, p2 *Position, oldRange int, newRange int, towers [][]*Tower, max *Position) (removeTowers []*Tower, addTowers []*Tower, unChangeTowers []*Tower) {
	start1, end1 := t.getPosLimit(p1, oldRange, max)
	start2, end2 := t.getPosLimit(p2, newRange, max)

	for x := start1.X; x <= end1.X; x++ {
		for y := start1.Y; y <= end1.Y; y++ {
			if isInRect(&Position{x, y}, start2, end2) {
				if unChangeTowers == nil {
					unChangeTowers = make([]*Tower, 1)
				}

				unChangeTowers = append(unChangeTowers, towers[x][y])
			} else {
				if removeTowers == nil {
					removeTowers = make([]*Tower, 1)
				}

				removeTowers = append(removeTowers, towers[x][y])
			}
		}
	}

	for x := start2.X; x <= end2.X; x++ {
		for y := start2.Y; y <= end2.Y; y++ {
			if !isInRect(&Position{x, y}, start1, end1) {
				if addTowers == nil {
					addTowers = make([]*Tower, 1)
				}

				addTowers = append(addTowers, towers[x][y])
			}
		}
	}
	return
}

// Get the postion limit of given range
// @param pos {Object} The center position
// @param range {Number} The range
// @param max {max} The limit, the result will not exceed the limit
// @return The pos limitition  
func (t *TowerAOI) getPosLimit(pos *Position, r int, max *Position) (start *Position, end *Position) {
	if start == nil {
		start = &Position{}
	}
	if end == nil {
		end = &Position{}
	}

	if pos.X-r < 0 {
		start.X = 0
		end.X = 2 * r
	} else if pos.X+r > max.X {
		end.X = max.X
		start.X = max.X - 2*r
	} else {
		start.X = pos.X - r
		end.X = pos.X + r
	}

	if pos.Y-r < 0 {
		start.Y = 0
		end.Y = 2 * r
	} else if pos.Y+r > max.Y {
		end.Y = max.Y
		start.Y = max.Y - 2*r
	} else {
		start.Y = pos.Y - r
		end.Y = pos.Y + r
	}

	if start.X < 0 {
		start.X = 0
	}
	if end.X > max.X {
		end.X = max.X
	}

	if start.Y < 0 {
		start.Y = 0
	}
	if end.Y > max.Y {
		end.Y = max.Y
	}
	return
}

// Check if the pos is in the rect
func isInRect(pos *Position, start *Position, end *Position) bool {
	return (pos.X >= start.X && pos.X <= end.X && pos.Y >= start.Y && pos.Y <= end.Y)
}

func (t *TowerAOI) GetWatchers(pos *Position, types []int) map[int]map[int64]int64 {
	if t.checkPos(pos) {
		p := t.transPos(pos)
		return t.towers[p.X][p.Y].GetWatchers(types)
	}
	return nil
}

// Combine map to arr
// @param arr {Array} The array to add the map to
// @param map {Map} The map to add to array
func addMap(result []int64, m map[int64]int64) []int64 {
	r := make([]int64, len(m))
	i := 0
	for _, v := range m {
		r[i] = v
		i++
	}
	return append(result, r...)
}

func (t *TowerAOI) addMapByTypes(result map[int]map[int64]int64, m map[int]map[int64]int64, types []int) map[int]map[int64]int64 {
	for i := 0; i < len(types); i++ {
		objType := types[i]

		_, r1 := m[objType]
		if !r1 {
			continue
		}

		_, r2 := result[objType]
		if !r2 {
			result[objType] = make(map[int64]int64)
		}

		for k, v := range m[objType] {
			result[objType][k] = v
		}
	}
	return result
}

// Create TowerAOI object
func CreateTowerAOI(config *Config) *TowerAOI {
	towerAOI := &TowerAOI{}
	towerAOI.Init(config)

	return towerAOI
}
