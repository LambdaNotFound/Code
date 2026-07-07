package design

/**
 * Design Excel
 *    OO design: Sheet, Cell
 *
 *    a cell w/ equation e.g. 1 + 1
 *    a cell reference other cells
 */

type Cell interface {
    SetContent(content string)
    Evaluate() string
    AddListener(c cell)
    RemoveListener(c cell)
}

type cell struct {
    content   string
    alias     string // 'A2'
    listeners []cell
}

/**
 * Small immutable structs → value receiver.
 *
 * Large or mutable structs → pointer receiver.
 */

func (c* cell) SetContent(content string) {
    c.content = content
}

type Sheet interface {
    UpdateCell(alias string, content string)

    updateListeners(alias string)
}

type sheet struct {
    cellsMap map[string]cell
}
