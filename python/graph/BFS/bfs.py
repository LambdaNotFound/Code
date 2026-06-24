from collections import deque


def bfs(start, graph):
    """Generic BFS over an adjacency structure."""
    queue = deque([start])
    visited = {start}  # mark on ENQUEUE, not dequeue, to avoid dupes

    while queue:
        node = queue.popleft()
        for neighbor in graph[node]:
            if neighbor not in visited:
                visited.add(neighbor)
                queue.append(neighbor)
    return visited


def bfs_levels(start, target, graph):
    queue = deque([start])
    visited = {start}
    steps = 0

    while queue:
        for _ in range(len(queue)):  # drain exactly one level
            node = queue.popleft()
            if node == target:
                return steps
            for neighbor in graph[node]:
                if neighbor not in visited:
                    visited.add(neighbor)
                    queue.append(neighbor)
        steps += 1
    return -1


def bfs_grid(grid, start, obstacle="#"):
    rows, cols = len(grid), len(grid[0])
    queue = deque([start])
    visited = {start}
    DIRS = [(-1, 0), (1, 0), (0, -1), (0, 1)]  # 4-directional

    while queue:
        r, c = queue.popleft()
        for dr, dc in DIRS:
            nr, nc = r + dr, c + dc
            if (
                0 <= nr < rows
                and 0 <= nc < cols
                and (nr, nc) not in visited
                and grid[nr][nc] != obstacle
            ):
                visited.add((nr, nc))
                queue.append((nr, nc))
    return visited


"""
Given:
  - computer_connections: a list of pairs (c1, c2) that represents a
  directional connection from c1 to c2 (not the otherway around),
  - infected_time: an integer that represent the seconds it take for one
  computer to infect its immediate neighbor
  - start_computer: an integer that represent the computer where it starts the infection
  
Return the minimum time it takes for every computer to be infected, if that is not 
possible, return -1.

Know that the start_computer is immediately infected so it takes Øs to infected it.

Example 1:
computer_connections: [(0,1), (0,2), (1,3)]
infected_time = 120 
start_computer = 0
      |---[1]---[3]
      |
[0]---|---[2]

It takes 120s for computer 0 to infect computer 1 and 2. It takes another
120s for computer 1 to infect computer 3. Thus, it takes 240s to infect all
computers.

Example 2:
computer_connections: [(0,1), (0,2)]
infected_time = 120
start_computer = 3 # not in the connections
=> return -1
"""


def time_to_infect(computer_connections, infected_time, start_computer):
    from collections import defaultdict, deque

    all_nodes = set()
    graph = defaultdict(list)
    for c1, c2 in computer_connections:
        graph[c1].append(c2)
        all_nodes.add(c1)
        all_nodes.add(c2)

    if start_computer not in all_nodes:
        return -1  # or handle the empty-connections case

    queue = deque([start_computer])
    visited = {start_computer}
    levels = 0

    while queue:
        for _ in range(len(queue)):
            node = queue.popleft()
            for neighbor in graph[node]:
                if neighbor not in visited:
                    visited.add(neighbor)
                    queue.append(neighbor)
        levels += 1  # note: this overcounts by 1 — levels ends at depth+1

    if visited != all_nodes:
        return -1
    return (levels - 1) * infected_time
