//! Data-structure and algorithm solutions, one module per leetcode.com tag,
//! mirroring `go/dsa/`. Each module is self-contained and defines its own
//! `Solution`, the way LeetCode's Rust template does.

// A solution file often shares its folder's name, as `go/dsa/heap/heap.go`
// does, which makes `heap/heap.rs` a module inside a module of the same
// name. Clippy warns on that by default; the mirror is the point here.
#![allow(clippy::module_inception)]

pub mod array;
pub mod backtracking;
pub mod bfs;
pub mod binary_search;
pub mod binary_tree;
pub mod bit_manipulation;
pub mod bitmask;
pub mod bst;
pub mod counting;
pub mod design;
pub mod dfs;
pub mod divide_and_conquer;
pub mod dp;
pub mod graph;
pub mod greedy;
pub mod hash_table;
pub mod heap;
pub mod linked_list;
pub mod math;
pub mod matrix;
pub mod memoization;
pub mod mono_queue;
pub mod mono_stack;
pub mod ordered_set;
pub mod prefix_sum;
pub mod shortest_path;
pub mod simulation;
pub mod sliding_window;
pub mod stack;
pub mod string;
pub mod topological_sort;
pub mod trie;
pub mod two_pointers;
pub mod types;
pub mod union_find;
