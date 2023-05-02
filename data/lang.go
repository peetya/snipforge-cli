package data

import "strings"

type FormatSnippet func(snippet string) string

type Language struct {
	Names             []string
	PreferredFileName string
	Format            FormatSnippet
}

var Languages = []Language{
	{Names: []string{"PHP", "Symfony", "Laravel"}, PreferredFileName: "snippet.php", Format: func(snippet string) string {
		// Sometimes PHP opening tag is missing from snippets from OpenAI API
		if !strings.HasPrefix(snippet, "<?php") {
			snippet = "<?php\n\n" + snippet
		}

		return snippet
	}},
	{Names: []string{"Python", "Django"}, PreferredFileName: "snippet.py"},
	{Names: []string{"Go", "Golang"}, PreferredFileName: "snippet.go"},
	{Names: []string{"JavaScript", "JS", "Node", "NodeJS", "Express"}, PreferredFileName: "snippet.js"},
	{Names: []string{"React"}, PreferredFileName: "snippet.jsx"},
	{Names: []string{"TypeScript", "TS"}, PreferredFileName: "snippet.ts"},
	{Names: []string{"Java"}, PreferredFileName: "Snippet.java"},
	{Names: []string{"C#", "Csharp"}, PreferredFileName: "snippet.cs"},
	{Names: []string{"C++", "Cplusplus"}, PreferredFileName: "snippet.cpp"},
	{Names: []string{"C"}, PreferredFileName: "snippet.c"},
	{Names: []string{"Ruby"}, PreferredFileName: "snippet.rb"},
	{Names: []string{"Rust"}, PreferredFileName: "snippet.rs"},
	{Names: []string{"Swift"}, PreferredFileName: "snippet.swift"},
	{Names: []string{"Kotlin"}, PreferredFileName: "snippet.kt"},
	{Names: []string{"Scala"}, PreferredFileName: "snippet.scala"},
	{Names: []string{"R"}, PreferredFileName: "snippet.r"},
	{Names: []string{"Dart"}, PreferredFileName: "snippet.dart"},
	{Names: []string{"Haskell"}, PreferredFileName: "snippet.hs"},
	{Names: []string{"Julia"}, PreferredFileName: "snippet.jl"},
	{Names: []string{"Perl"}, PreferredFileName: "snippet.pl"},
	{Names: []string{"Lua"}, PreferredFileName: "snippet.lua"},
	{Names: []string{"Elixir"}, PreferredFileName: "snippet.ex"},
	{Names: []string{"Clojure"}, PreferredFileName: "snippet.clj"},
	{Names: []string{"Erlang"}, PreferredFileName: "snippet.erl"},
	{Names: []string{"OCaml"}, PreferredFileName: "snippet.ml"},
	{Names: []string{"Racket"}, PreferredFileName: "snippet.rkt"},
	{Names: []string{"F#"}, PreferredFileName: "snippet.fs"},
	{Names: []string{"Bash", "Shell", "sh"}, PreferredFileName: "snippet.sh"},
	{Names: []string{"PowerShell"}, PreferredFileName: "snippet.ps1"},
	{Names: []string{"SQL"}, PreferredFileName: "snippet.sql"},
	{Names: []string{"HTML"}, PreferredFileName: "snippet.html"},
	{Names: []string{"CSS"}, PreferredFileName: "snippet.css"},
	{Names: []string{"JSON"}, PreferredFileName: "snippet.json"},
	{Names: []string{"YAML"}, PreferredFileName: "snippet.yaml"},
	{Names: []string{"XML"}, PreferredFileName: "snippet.xml"},
	{Names: []string{"Markdown", "MD"}, PreferredFileName: "snippet.md"},
	{Names: []string{"Dockerfile"}, PreferredFileName: "Dockerfile"},
	{Names: []string{"Makefile"}, PreferredFileName: "Makefile"},
	{Names: []string{"INI"}, PreferredFileName: "snippet.ini"},
	{Names: []string{"TOML"}, PreferredFileName: "snippet.toml"},
	{Names: []string{"CSV"}, PreferredFileName: "snippet.csv"},
	{Names: []string{"Terraform"}, PreferredFileName: "snippet.tf"},
	{Names: []string{"HCL"}, PreferredFileName: "snippet.hcl"},
	{Names: []string{"GraphQL"}, PreferredFileName: "snippet.gql"},
	{Names: []string{"HTTP"}, PreferredFileName: "snippet.http"},
	{Names: []string{"Shell"}, PreferredFileName: "snippet.sh"},
}
