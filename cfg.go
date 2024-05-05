package main

import (
    "os"
    "gopkg.in/yaml.v2"
    "github.com/google/uuid"
)

type Config struct {
    Source string `yaml:"source"`
    Listen struct {
      Service string `yaml:"service"`
      Metrics string `yaml:"metrics"`
    } `yaml:"http"`
    Log struct {
      File string `yaml:"file"`
      Request string `yaml:"request"`
    } `yaml:"logger"`
    Sys struct {
      Threads int `yaml:"threads"`
      Maxout int `yaml:"maxout"`
    } `yaml:"system"`
    Uuid string
    Host string
}


var CFG Config

func Config_init(fn string) error {
  f, err := os.Open(fn)
  if err != nil {
    return err
  }
  defer f.Close()

  var cfg Config
  decoder := yaml.NewDecoder(f)
  err = decoder.Decode(&cfg)
  if err != nil {
    return err
  }
  CFG=cfg
  CFG.Uuid=uuid.New().String()
  CFG.Host, _ =os.Hostname()
  return nil
}
