package main

import (
	"github.com/deissh/osu-lazer/server/cmd/risuserver/commands"
	"os"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}

//config := model.SqlSettings{}
//config.SetDefaults(true)
//
//testStore = store.NewTimerLayer(
//sqlstore.NewSqlSupplier(config, nil),
//nil,
//)
//
//testStore.User().GetAll()
