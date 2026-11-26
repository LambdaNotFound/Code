package interview

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

/**
 * FileSystem — flat baseline file store: a single map[name]size with no
 * ownership, TTL, or history. AddFile/GetFileSize/DeleteFile are O(1);
 * ListFiles is O(n log n), sorting by size desc then name asc.
 */
type FileSystem struct {
	files map[string]int // name -> size
}

func NewFileSystem() *FileSystem {
	return &FileSystem{files: make(map[string]int)}
}

func (fs *FileSystem) AddFile(name string, size int) bool {
	if _, exists := fs.files[name]; exists {
		return false
	}
	fs.files[name] = size
	return true
}

func (fs *FileSystem) GetFileSize(name string) *int {
	size, ok := fs.files[name]
	if !ok {
		return nil
	}
	return &size
}

func (fs *FileSystem) DeleteFile(name string) bool {
	if _, ok := fs.files[name]; !ok {
		return false
	}
	delete(fs.files, name)
	return true
}

func (fs *FileSystem) ListFiles() []string {
	type entry struct {
		name string
		size int
	}
	all := make([]entry, 0, len(fs.files))
	for name, size := range fs.files {
		all = append(all, entry{name, size})
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].size != all[j].size {
			return all[i].size > all[j].size
		}
		return all[i].name < all[j].name
	})
	result := make([]string, len(all))
	for i, e := range all {
		result[i] = e.name + "(" + strconv.Itoa(e.size) + ")"
	}
	return result
}

// PROGRESSIVE

type FileV1 struct {
	name         string
	size         int
	originalSize int // for compression (Problem 4)
	compressed   bool
	ownerId      string
}

type UserV1 struct {
	id       string
	capacity int
	files    map[string]*FileV1 // name -> file
}

/**
 * ProgressiveFS — progressive mock-interview file system, levels 1-4:
 *   level 1: AddFile/GetFileSize/DeleteFile on a flat file map
 *   level 2: FindFiles by prefix+suffix (via shared formatAndSort helper)
 *   level 3: per-user ownership with capacity limits (AddUser,
 *            AddFileByUser tracks usage by summing owned file sizes)
 *   level 4: UpdateCapacity evicts largest-file-first (name-desc tie-break)
 *            until under the new cap; CopyFile respects owner capacity;
 *            Compress/DecompressFile halve/restore size via originalSize
 *
 * Storage: files map[name]*FileV1 for O(1) lookup, cross-referenced by
 * users[id].files map[name]*FileV1 (same pointers) so per-user usage can
 * be recomputed by summing that user's file map without touching fs.files.
 *
 * Complexity: AddFile/GetFileSize/DeleteFile/CopyFile/Compress/Decompress
 * are O(1) (O(u) for AddFileByUser/CopyFile's capacity scan, u = user's
 * file count). FindFiles is O(n log n). UpdateCapacity is O(u^2) worst
 * case (linear victim scan repeated per eviction).
 */
type ProgressiveFS struct {
	files map[string]*FileV1
	users map[string]*UserV1
}

func NewProgressiveFS() *ProgressiveFS {
	return &ProgressiveFS{
		files: make(map[string]*FileV1),
		users: make(map[string]*UserV1),
	}
}

// ==================== LEVEL 1 ====================

func (fs *ProgressiveFS) AddFile(name string, size int) bool {
	if _, exists := fs.files[name]; exists {
		return false
	}
	fs.files[name] = &FileV1{name: name, size: size, originalSize: size}
	return true
}

func (fs *ProgressiveFS) GetFileSize(name string) *int {
	f, ok := fs.files[name]
	if !ok {
		return nil
	}
	return &f.size
}

func (fs *ProgressiveFS) DeleteFile(name string) bool {
	f, ok := fs.files[name]
	if !ok {
		return false
	}
	if f.ownerId != "" {
		if u, ok := fs.users[f.ownerId]; ok {
			delete(u.files, name)
		}
	}
	delete(fs.files, name)
	return true
}

// ==================== LEVEL 2 ====================

func (fs *ProgressiveFS) FindFiles(prefix, suffix string) []string {
	var matched []*FileV1
	for _, f := range fs.files {
		if strings.HasPrefix(f.name, prefix) && strings.HasSuffix(f.name, suffix) {
			matched = append(matched, f)
		}
	}
	return formatAndSortFiles(matched)
}

func formatAndSortFiles(files []*FileV1) []string {
	sort.Slice(files, func(i, j int) bool {
		if files[i].size != files[j].size {
			return files[i].size > files[j].size
		}
		return files[i].name < files[j].name
	})
	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.name + "(" + strconv.Itoa(f.size) + ")"
	}
	return result
}

// ==================== LEVEL 3 ====================

func (fs *ProgressiveFS) AddUser(userId string, capacity int) bool {
	if _, exists := fs.users[userId]; exists {
		return false
	}
	fs.users[userId] = &UserV1{
		id:       userId,
		capacity: capacity,
		files:    make(map[string]*FileV1),
	}
	return true
}

func (fs *ProgressiveFS) AddFileByUser(userId, name string, size int) bool {
	if _, exists := fs.files[name]; exists {
		return false
	}
	u, ok := fs.users[userId]
	if !ok {
		return false
	}
	// Check capacity
	used := 0
	for _, f := range u.files {
		used += f.size
	}
	if used+size > u.capacity {
		return false
	}
	f := &FileV1{name: name, size: size, originalSize: size, ownerId: userId}
	fs.files[name] = f
	u.files[name] = f
	return true
}

// ==================== LEVEL 4 ====================

func (fs *ProgressiveFS) UpdateCapacity(userId string, newCapacity int) *int {
	u, ok := fs.users[userId]
	if !ok {
		return nil
	}
	u.capacity = newCapacity

	// Calculate current usage
	used := 0
	for _, f := range u.files {
		used += f.size
	}

	removed := 0
	for used > newCapacity {
		// Find largest file, tie-break by name descending (remove "last" alphabetically)
		var victim *FileV1
		for _, f := range u.files {
			if victim == nil || f.size > victim.size ||
				(f.size == victim.size && f.name > victim.name) {
				victim = f
			}
		}
		used -= victim.size
		delete(fs.files, victim.name)
		delete(u.files, victim.name)
		removed++
	}

	return &removed
}

func (fs *ProgressiveFS) CopyFile(fromName, toName string) bool {
	src, ok := fs.files[fromName]
	if !ok {
		return false
	}
	if _, exists := fs.files[toName]; exists {
		return false
	}
	// Check owner capacity if owned
	if src.ownerId != "" {
		u := fs.users[src.ownerId]
		used := 0
		for _, f := range u.files {
			used += f.size
		}
		if used+src.size > u.capacity {
			return false
		}
	}
	f := &FileV1{
		name:         toName,
		size:         src.size,
		originalSize: src.originalSize,
		ownerId:      src.ownerId,
	}
	fs.files[toName] = f
	if src.ownerId != "" {
		fs.users[src.ownerId].files[toName] = f
	}
	return true
}

func (fs *ProgressiveFS) CompressFile(name string) bool {
	f, ok := fs.files[name]
	if !ok || f.compressed {
		return false
	}
	f.size = f.size / 2
	f.compressed = true
	return true
}

func (fs *ProgressiveFS) DecompressFile(name string) bool {
	f, ok := fs.files[name]
	if !ok || !f.compressed {
		return false
	}
	f.size = f.originalSize
	f.compressed = false
	return true
}

type FileV2 struct {
	name     string
	size     int
	expireAt int // 0 = never expires
}

type FileSnapshot struct {
	timestamp int
	files     map[string]*FileV2 // deep copy of state at that time
}

/**
 * FileServer — separate progressive problem: a timestamped file server,
 * levels 1-4:
 *   level 1: FileUpload/FileGet/FileCopy (thin wrappers with timestamp=0)
 *   level 2: FileSearch by prefix, capped to top 10 by size desc/name asc
 *   level 3: *At variants take an explicit timestamp; uploads carry a TTL
 *            and files are treated as absent once timestamp >= expireAt
 *            (panics on upload-name conflict or copy of a dead/missing
 *            source, rather than returning a bool/error)
 *   level 4: saveSnapshot deep-copies the whole file map after every
 *            mutation; Rollback finds the latest snapshot at or before a
 *            timestamp and restores it, shifting each restored file's
 *            expireAt forward by the elapsed time so TTLs measured from
 *            the rollback point behave as if they'd been ticking all along
 *
 * Complexity: FileUpload/FileGet/FileCopy are O(1) plus O(n) for the
 * post-mutation snapshot copy (n = live file count). FileSearch is
 * O(n log n) for the sort. Rollback is O(s + n) (s = snapshot count
 * scanned back-to-front, n = files restored).
 */
type FileServer struct {
	files     map[string]*FileV2
	snapshots []FileSnapshot // ordered by timestamp ascending
}

func NewFileServer() *FileServer {
	return &FileServer{
		files: make(map[string]*FileV2),
	}
}

// ==================== LEVEL 1 ====================

func (fs *FileServer) FileUpload(fileName string, size int) {
	fs.FileUploadAt(0, fileName, size, 0)
}

func (fs *FileServer) FileGet(fileName string) *int {
	return fs.FileGetAt(0, fileName)
}

func (fs *FileServer) FileCopy(source, dest string) {
	fs.FileCopyAt(0, source, dest)
}

// ==================== LEVEL 2 ====================

func (fs *FileServer) FileSearch(prefix string) []string {
	return fs.FileSearchAt(0, prefix)
}

// ==================== LEVEL 3 ====================

func (fs *FileServer) FileUploadAt(timestamp int, fileName string, size int, ttl int) {
	if f, exists := fs.files[fileName]; exists {
		// only conflict if file is still alive
		if f.expireAt == 0 || timestamp < f.expireAt {
			panic(fmt.Sprintf("file %s already exists", fileName))
		}
	}
	f := &FileV2{name: fileName, size: size}
	if ttl > 0 {
		f.expireAt = timestamp + ttl
	}
	fs.files[fileName] = f
	fs.saveSnapshot(timestamp)
}

func (fs *FileServer) FileGetAt(timestamp int, fileName string) *int {
	f, ok := fs.files[fileName]
	if !ok {
		return nil
	}
	if f.expireAt != 0 && timestamp >= f.expireAt {
		return nil
	}
	return &f.size
}

func (fs *FileServer) FileCopyAt(timestamp int, source, dest string) {
	src, ok := fs.files[source]
	if !ok || (src.expireAt != 0 && timestamp >= src.expireAt) {
		panic(fmt.Sprintf("source file %s does not exist", source))
	}
	// Copy inherits no TTL — fresh file at dest
	fs.files[dest] = &FileV2{name: dest, size: src.size}
	fs.saveSnapshot(timestamp)
}

func (fs *FileServer) FileSearchAt(timestamp int, prefix string) []string {
	var matched []*FileV2
	for _, f := range fs.files {
		if f.expireAt != 0 && timestamp >= f.expireAt {
			continue
		}
		if strings.HasPrefix(f.name, prefix) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].size != matched[j].size {
			return matched[i].size > matched[j].size // desc size
		}
		return matched[i].name < matched[j].name // asc name tie-break
	})
	if len(matched) > 10 {
		matched = matched[:10]
	}
	result := make([]string, len(matched))
	for i, f := range matched {
		result[i] = f.name
	}
	return result
}

// ==================== LEVEL 4 ====================

func (fs *FileServer) saveSnapshot(timestamp int) {
	snapshot := make(map[string]*FileV2)
	for k, f := range fs.files {
		// deep copy
		copy := *f
		snapshot[k] = &copy
	}
	fs.snapshots = append(fs.snapshots, FileSnapshot{
		timestamp: timestamp,
		files:     snapshot,
	})
}

func (fs *FileServer) Rollback(timestamp int) {
	// find latest snapshot at or before timestamp
	var target *FileSnapshot
	for i := len(fs.snapshots) - 1; i >= 0; i-- {
		if fs.snapshots[i].timestamp <= timestamp {
			target = &fs.snapshots[i]
			break
		}
	}
	if target == nil {
		return
	}

	// Restore, recalculating expireAt
	fs.files = make(map[string]*FileV2)
	elapsed := timestamp - target.timestamp
	for k, f := range target.files {
		restored := *f
		if restored.expireAt != 0 {
			restored.expireAt += elapsed // shift expiry forward by elapsed time
		}
		fs.files[k] = &restored
	}
}
