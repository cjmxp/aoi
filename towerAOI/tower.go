package aoi

type Tower struct {
	ids      map[int64]int64
	watchers map[int]map[int64]int64
	typeMap  map[int]map[int64]int64
	size     int
}

type Object struct {
	ID      int64
	ObjType int
}

func (t *Tower) Init() {
	t.ids = make(map[int64]int64)
	t.watchers = make(map[int]map[int64]int64)
	t.typeMap = make(map[int]map[int64]int64)
}

// Add an object to tower
func (t *Tower) Add(obj *Object) bool {
	id := obj.ID
	objType := obj.ObjType

	if objType != 0 {
		_, ok := t.typeMap[objType]
		if !ok {
			t.typeMap[objType] = make(map[int64]int64)
		}
		if t.typeMap[objType][id] == id {
			return false
		}

		t.ids[id] = id
		t.typeMap[objType][id] = id
		t.size++
		return true
	}
	return false
}

// Remove an object from this tower
func (t *Tower) Remove(obj *Object) {
	id := obj.ID
	objType := obj.ObjType

	_, ok := t.ids[id]
	if ok {
		delete(t.ids, id)
	}

	if objType != 0 {
		_, r := t.typeMap[objType][id]
		if r {
			delete(t.typeMap[objType], id)
		}
	}
	t.size--
}

// Get all object ids in this tower
func (t *Tower) GetIds() map[int64]int64 {
	return t.ids
}

// Add watcher to tower
func (t *Tower) AddWatcher(watcher *Object) {
	id := watcher.ID
	objType := watcher.ObjType

	if objType != 0 {
		_, ok := t.watchers[objType]
		if !ok {
			t.watchers[objType] = make(map[int64]int64)
		}
		t.watchers[objType][id] = id
	}
}

// Remove watcher from tower
func (t *Tower) RemoveWatcher(watcher *Object) {
	id := watcher.ID
	objType := watcher.ObjType

	if objType != 0 {
		_, ok := t.watchers[objType][id]
		if ok {
			delete(t.watchers[objType], id)
		}
	}
}

// Get all watchers by the given types in this tower
func (t *Tower) GetWatchers(types []int) map[int]map[int64]int64 {
	result := make(map[int]map[int64]int64)

	if types != nil && len(types) > 0 {
		for i := 0; i < len(types); i++ {
			objType := types[i]
			_, ok := t.watchers[objType]
			if ok {
				result[objType] = t.watchers[objType]
			}
		}
	}
	return result
}

// Get object ids of given types in this tower
func (t *Tower) GetIdsByTypes(types []int) map[int]map[int64]int64 {
	result := make(map[int]map[int64]int64)

	for i := 0; i < len(types); i++ {
		objType := types[i]
		_, ok := t.typeMap[objType]

		if ok {
			result[objType] = t.typeMap[objType]
		}
	}

	return result
}

// Create Tower object
func CreateTower() (tower *Tower) {
	tower = &Tower{}
	tower.Init()

	return
}
