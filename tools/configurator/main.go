package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/wolfeidau/unflatten"
)

var (
	oldConfigFile   string
	newConfigFile   string
	finalConfigFile string
)

func init() {
	flag.StringVar(&oldConfigFile, "old-config", "tikv.toml", "old config file name")
	flag.StringVar(&newConfigFile, "new-config", "new-tikv.toml", "new config file name")
	flag.StringVar(&finalConfigFile, "final-config", "final-tikv.toml", "final config file name")
	flag.Parse()
}

func main() {
	var oldConfig, newConfig, finalConfig map[string]interface{}
	if _, err := toml.DecodeFile(oldConfigFile, &oldConfig); err != nil {
		log.Fatalf("failed to parse old config file %s: %v", oldConfigFile, err)
	}
	if _, err := toml.DecodeFile(newConfigFile, &newConfig); err != nil {
		log.Fatalf("failed to parse new config file %s: %v", newConfigFile, err)
	}

	log.Printf("old config: %v", oldConfig)
	log.Printf("new config: %v", newConfig)

	oldFlattened := unflatten.Flatten(oldConfig, unflatten.JoinWithDot)
	newFlattened := unflatten.Flatten(newConfig, unflatten.JoinWithDot)

	finalFlattened := map[string]interface{}{}
	for key := range newFlattened {
		if _, ok := oldFlattened[key]; !ok {
			finalFlattened[key] = newFlattened[key]
		}
		finalFlattened[key] = oldFlattened[key]
	}

	finalConfig = unflatten.Unflatten(finalFlattened, unflatten.SplitByDot)
	log.Printf("final config: %v", finalConfig)
	file, err := os.OpenFile(finalConfigFile, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		log.Fatalf("failed to create or open final config file %s: %v", finalConfigFile, err)
	}
	enc := toml.NewEncoder(file)
	if err := enc.Encode(finalConfig); err != nil {
		log.Fatalf("failed to write final config %v to file %s: %v", finalConfig, finalConfigFile, err)
	}
}
