package cmg

import (
	"encoding/json"
	"fmt"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/mapgen"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"log"
	"net/http"

	"go.opencensus.io/plugin/ochttp"
)


func getMap(w http.ResponseWriter, r *http.Request) {
	rules := game.GameRules{
		GameType:             0,
		MinimumScore:         165,
		MaximumScore:         361,
		MaxOver300:           14,
		MaximumResourceScore: 130,
		MinimumResourceScore: 30,
	}

	gameType := game.NormalGame
	if rules.GameType == 0 {
		gameType = game.NormalGame
	} else if rules.GameType == 1 {
		gameType = game.LargeGame
	}
	verbose := false
	numberOfLoops := 1
	totalGenerations := 0

	board := mapgen.MapGenerationAttempt(gameType, verbose)
	for i := 0; i < numberOfLoops; i++ {
		for !board.IsValid(rules, gameType, verbose) {
			totalGenerations++
			if totalGenerations > 1501 {
				log.Fatal("Can not generate a map... (1000+ runs)")
			}
			board = mapgen.MapGenerationAttempt(gameType, verbose)
		}
	}
	var content = model.Map{
		GameType: "Normal",
		Board:    board.Board,
	}
	b, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintln(w, string(b))
}

func Cmg(w http.ResponseWriter, r *http.Request) {
	traced := &ochttp.Handler{
		Handler: http.HandlerFunc(getMap),
	}
	traced.ServeHTTP(w, r)
}
