package main

func main() {
    _, err := NewStore("./store.db")
    if err != nil {
        panic(err)
    }
}
