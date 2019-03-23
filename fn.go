package cmg

import (
	"encoding/json"
	"fmt"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/mapgen"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"log"
	"net/http"
	"strconv"

	"go.opencensus.io/plugin/ochttp"
)

func extractIntParamOrDefault(request *http.Request, paramName string, defaultValue int) int {
	paramValues :=  request.URL.Query().Get(paramName)
	if len(paramValues) <= 0 {
		return defaultValue
	}
	paramValue := fmt.Sprintf("%v", paramValues)
	log.Println(fmt.Sprintf("Param check, name: %s, value: %v, values: %v", paramName, paramValue, paramValues))
	intValue, err := strconv.Atoi(string(paramValue))
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getMap(w http.ResponseWriter, r *http.Request) {
	log.Println("Request params were:", r.URL.Query())

	gameTypeValue := 0
	gameTypeParam := r.URL.Query().Get("type")
	gameTypeParamValue := fmt.Sprintf("%v",gameTypeParam)
	if gameTypeParamValue == "large" {
		gameTypeValue = 1
	}

	min := extractIntParamOrDefault(r, "min", 165)
	max := extractIntParamOrDefault(r, "max", 361)
	max300 := extractIntParamOrDefault(r, "max300", 14)
	maxr := extractIntParamOrDefault(r, "maxr", 130)
	minr := extractIntParamOrDefault(r, "minr", 30)

	rules := game.GameRules{
		GameType:             gameTypeValue,
		MinimumScore:         min,
		MaximumScore:         max,
		MaxOver300:           max300,
		MaximumResourceScore: maxr,
		MinimumResourceScore: minr,
	}
	log.Println("Rules: ", rules)

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
		GameType: gameType.Name,
		Board:    board.Board,
	}
	b, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(b))
}



func Cmg(w http.ResponseWriter, r *http.Request) {
	traced := &ochttp.Handler{
		Handler: http.HandlerFunc(getMap),
	}
	traced.ServeHTTP(w, r)
}
