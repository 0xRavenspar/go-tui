package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rivo/tview"
)

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"github.com/rivo/tview"
// )

// define an item structure that will hold the stock information
type Item struct {
	Name  string `json:"anme"`
	Stock int    `json:"stock"`
}

var (
	inventory     = []Item{}
	inventoryFile = "inventory.json"
)

// load inventory func
func loadInventory() {
	if _, err := os.Stat(inventoryFile); err == nil {
		{
			data, err := os.ReadFile(inventoryFile)
			if err != nil {
				log.Fatal("error reading inventory file! - ", err)
			}
			json.Unmarshal(data, &inventory)
		}
	}
}

// save inventory func
func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", "")
	if err != nil {
		log.Fatal("Error saving inventory! - ", err)
	}
	os.WriteFile(inventoryFile, data, 0644)
}

// delete item function
func deleteItem(index int) {
	if index < 0 || index >= len(inventory) {
		fmt.Println("Invalid item index!")
		return
	}
	inventory = append(inventory[:index], inventory[index+1:]...)
	saveInventory()
}

func main() {
	//Create a new TUI app

	app := tview.NewApplication()
	loadInventory()
	inventoryList := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)

	inventoryList.SetBorder(true).SetTitle("Inventory items")

	refreshInventory := func() {
		inventoryList.Clear()
		if len(inventory) == 0 {
			fmt.Fprintln(inventoryList, "No items in inventory")
		} else {
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Stock: %d)\n", i+1, item.Name, item.Stock)
			}
		}
	}

	//Creating three input fields

	itemNameInput := tview.NewInputField().SetLabel("item Name: ")
	itemStockInput := tview.NewInputField().SetLabel("STock: ")
	itemIDInput := tview.NewInputField().SetLabel("Item ID to delete: ")

	form := tview.NewForm().
		AddFormItem(itemNameInput).
		AddFormItem(itemStockInput).
		AddFormItem(itemIDInput).
		AddButton("Add Item", func() {
			name := itemNameInput.GetText()
			stock := itemStockInput.GetText()
			if name != "" && stock != "" {
				quantity, err := strconv.Atoi(stock)
				if err != nil {
					fmt.Fprintln(inventoryList, "Invalid stock value")
					return
				}
				inventory = append(inventory, Item{Name: name, Stock: quantity})
				saveInventory()
				refreshInventory()
				itemNameInput.SetText("")
				itemStockInput.SetText("")
			}
		}).
		AddButton("Delete item", func() {
			idStr := itemIDInput.GetText()
			if idStr == "" {
				fmt.Println(inventoryList, "Please enter an item ID to delete.")
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Fprintln(inventoryList, "Invalid item ID.")
				return
			}
			deleteItem(id - 1)
			fmt.Fprintf(inventoryList, "Item [%d] deleted.\n", id)
			refreshInventory()
			itemIDInput.SetText("")
		}).
		AddButton("Exit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("manga inventory").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(inventoryList, 0, 1, false).
		AddItem(form, 0, 1, true)

	refreshInventory()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
