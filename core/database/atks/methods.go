package atks

import (
	"fmt"
	"api/core/models/floods"
)

func (conn *Instance) LoadFloodMethods() error {
    rows, err := conn.conn.Query(`
      SELECT name, sname, type, description, subnet, vip
        FROM flood_methods
    `)
    if err != nil {
        return fmt.Errorf("LoadFloodMethods: %w", err)
    }
    defer rows.Close()

    // re‑init your in‑memory map
    floods.Methods = make(map[string]*floods.Method, 0)
    for rows.Next() {
        var m floods.Method
        if err := rows.Scan(
            &m.Name,
            &m.DisplayName,
            &m.Type,
            &m.Description,
            &m.Subnet,
            &m.VIP,
        ); err != nil {
            return fmt.Errorf("LoadFloodMethods: scan: %w", err)
        }
        floods.Methods[m.Name] = &m
    }
    if err := rows.Err(); err != nil {
        return fmt.Errorf("LoadFloodMethods: rows: %w", err)
    }
    return nil
}



func (conn *Instance) AddFloodMethod(m *floods.Method) error {
    _, err := conn.conn.Exec(`
        INSERT INTO flood_methods 
          (name, sname, type, description, subnet, vip)
        VALUES (?, ?, ?, ?, ?, ?)
    `, m.Name, m.DisplayName, m.Type, m.Description, m.Subnet, m.VIP)
    if err != nil {
        return fmt.Errorf("AddFloodMethod: %w", err)
    }
    // update in‑memory cache
    floods.Methods[m.Name] = m
    return nil
}

func (conn *Instance) UpdateFloodMethod(m *floods.Method, originalName string) error {
    _, err := conn.conn.Exec(`
        UPDATE flood_methods
           SET name        = ?,
               sname       = ?,
               type        = ?,
               description = ?,
               subnet      = ?,
               vip         = ?
         WHERE name        = ?
    `, m.Name, m.DisplayName, m.Type, m.Description, m.Subnet, m.VIP, originalName)
    if err != nil {
        return fmt.Errorf("UpdateFloodMethod: %w", err)
    }
    // keep your in‑memory map in sync:
    delete(floods.Methods, originalName)
    floods.Methods[m.Name] = m
    return nil
}

func (conn *Instance) DeleteFloodMethod(name string) error {
    _, err := conn.conn.Exec(`
        DELETE FROM flood_methods
         WHERE name = ?
    `, name)
    if err != nil {
        return fmt.Errorf("DeleteFloodMethod: %w", err)
    }
    // remove from in‑memory cache
    delete(floods.Methods, name)
    return nil
}