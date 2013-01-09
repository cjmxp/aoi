package aoi

import (
	"fmt"
	"testing"
)

func testTower(t *testing.T) {
	tower := CreateTower()

	obj1 := &Object{ID: 100, ObjType: 1}
	tower.Add(obj1)
	obj2 := &Object{ID: 101, ObjType: 1}
	tower.Add(obj2)
	obj3 := &Object{ID: 103, ObjType: 2}
	tower.Add(obj3)
	obj4 := &Object{ID: 104, ObjType: 3}
	tower.Add(obj4)

	fmt.Printf("%v\n", tower.GetIds())
	fmt.Printf("%v\n", tower.GetIdsByTypes([]int{2}))
	fmt.Printf("%v\n", tower.GetIdsByTypes([]int{1, 3}))

	tower.Remove(obj2)
	fmt.Printf("%v\n", tower.GetIdsByTypes([]int{1}))
	tower.Remove(obj1)
	fmt.Printf("%v\n", tower.GetIdsByTypes([]int{1}))
	fmt.Printf("%v\n", tower.GetIds())

	watcher1 := &Object{ID: 201, ObjType: 21}
	tower.AddWatcher(watcher1)
	watcher2 := &Object{ID: 201, ObjType: 22}
	tower.AddWatcher(watcher2)
	fmt.Printf("%v\n", tower.GetWatchers([]int{21, 22}))

	tower.RemoveWatcher(watcher2)
	fmt.Printf("%v\n", tower.GetWatchers([]int{21, 22}))

	fmt.Println(tower)
}

func testTowerAOI(t *testing.T) {
	config := &Config{
		M: &MapConfig{
			&Rectangle{
				10000,
				10000,
			},
			1001,
		},
		T: &TowerConfig{
			&Rectangle{
				100,
				100,
			},
		},
	}
	towerAOI := CreateTowerAOI(config)
}
