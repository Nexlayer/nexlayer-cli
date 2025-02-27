// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import (
	"path/filepath"
)

// TechStackDefinitions defines patterns for detecting technology stacks
var TechStackDefinitions = map[string]StackDefinition{
	"nextjs-supabase-langchain": {
		Name:        "Next.js + Supabase + LangChain",
		Description: "Modern stack for AI-powered applications using Next.js, Supabase, and LangChain",
		Components: Components{
			Frontend: []string{"nextjs"},
			Backend:  []string{"supabase"},
			Database: []string{"postgres", "pgvector"},
			AI:       []string{"langchain"},
		},
		RequiredComponents: []string{"nextjs", "supabase", "langchain"},
		OptionalComponents: []string{"pgvector", "tailwind", "stripe"},
		MainPatterns: []DetectionPattern{
			{
				Type:       PatternDependency,
				Pattern:    "next",
				Path:       "package.json",
				Confidence: 0.6,
			},
			{
				Type:       PatternDependency,
				Pattern:    "@supabase/supabase-js",
				Path:       "package.json",
				Confidence: 0.6,
			},
			{
				Type:       PatternDependency,
				Pattern:    "langchain",
				Path:       "package.json",
				Confidence: 0.6,
			},
		},
		ExtraPatterns: []DetectionPattern{
			{
				Type:       PatternFile,
				Pattern:    filepath.Join("app", "api", "(.+)", "route.ts"),
				Path:       "",
				Confidence: 0.1,
			},
			{
				Type:       PatternEnvironment,
				Pattern:    "SUPABASE_URL",
				Path:       ".env",
				Confidence: 0.1,
			},
			{
				Type:       PatternEnvironment,
				Pattern:    "SUPABASE_ANON_KEY",
				Path:       ".env",
				Confidence: 0.1,
			},
			{
				Type:       PatternContent,
				Pattern:    "import.*?createClient.*?supabase",
				Path:       filepath.Join("**", "*.{js,ts,jsx,tsx}"),
				Confidence: 0.1,
			},
			{
				Type:       PatternContent,
				Pattern:    "import.*from.*langchain",
				Path:       filepath.Join("**", "*.{js,ts,jsx,tsx}"),
				Confidence: 0.1,
			},
		},
	},

	"nextjs-supabase-openai": {
		Name:        "Next.js + Supabase + OpenAI",
		Description: "Modern stack for AI-powered applications using Next.js, Supabase, and OpenAI",
		Components: Components{
			Frontend: []string{"nextjs"},
			Backend:  []string{"supabase"},
			Database: []string{"postgres"},
			AI:       []string{"openai"},
		},
		RequiredComponents: []string{"nextjs", "supabase", "openai"},
		OptionalComponents: []string{"tailwind", "stripe"},
		MainPatterns: []DetectionPattern{
			{
				Type:       PatternDependency,
				Pattern:    "next",
				Path:       "package.json",
				Confidence: 0.6,
			},
			{
				Type:       PatternDependency,
				Pattern:    "@supabase/supabase-js",
				Path:       "package.json",
				Confidence: 0.6,
			},
			{
				Type:       PatternDependency,
				Pattern:    "openai",
				Path:       "package.json",
				Confidence: 0.6,
			},
		},
		ExtraPatterns: []DetectionPattern{
			{
				Type:       PatternEnvironment,
				Pattern:    "OPENAI_API_KEY",
				Path:       ".env",
				Confidence: 0.1,
			},
			{
				Type:       PatternContent,
				Pattern:    "import.*?OpenAI",
				Path:       filepath.Join("**", "*.{js,ts,jsx,tsx}"),
				Confidence: 0.1,
			},
		},
	},

	"django-react": {
		Name:        "Django + React",
		Description: "Full-stack web application using Django backend and React frontend",
		Components: Components{
			Frontend: []string{"react"},
			Backend:  []string{"django"},
			Database: []string{"postgres"},
		},
		RequiredComponents: []string{"django", "react"},
		OptionalComponents: []string{"postgres", "redis"},
		MainPatterns: []DetectionPattern{
			{
				Type:       PatternFile,
				Pattern:    "manage.py",
				Path:       "",
				Confidence: 0.5,
			},
			{
				Type:       PatternFile,
				Pattern:    filepath.Join("frontend", "package.json"),
				Path:       "",
				Confidence: 0.3,
			},
			{
				Type:       PatternDependency,
				Pattern:    "react",
				Path:       filepath.Join("frontend", "package.json"),
				Confidence: 0.3,
			},
		},
		ExtraPatterns: []DetectionPattern{
			{
				Type:       PatternFile,
				Pattern:    "requirements.txt",
				Path:       "",
				Confidence: 0.1,
			},
			{
				Type:       PatternContent,
				Pattern:    "django",
				Path:       "requirements.txt",
				Confidence: 0.1,
			},
			{
				Type:       PatternContent,
				Pattern:    "psycopg2",
				Path:       "requirements.txt",
				Confidence: 0.1,
			},
		},
	},

	"express-mongodb": {
		Name:        "Express + MongoDB",
		Description: "Backend API using Express.js with MongoDB database",
		Components: Components{
			Backend:  []string{"express", "node"},
			Database: []string{"mongodb"},
		},
		RequiredComponents: []string{"express", "mongodb"},
		OptionalComponents: []string{"mongoose"},
		MainPatterns: []DetectionPattern{
			{
				Type:       PatternDependency,
				Pattern:    "express",
				Path:       "package.json",
				Confidence: 0.5,
			},
			{
				Type:       PatternDependency,
				Pattern:    "mongodb",
				Path:       "package.json",
				Confidence: 0.4,
			},
		},
		ExtraPatterns: []DetectionPattern{
			{
				Type:       PatternEnvironment,
				Pattern:    "MONGO_URI",
				Path:       ".env",
				Confidence: 0.1,
			},
			{
				Type:       PatternContent,
				Pattern:    "mongoose.connect",
				Path:       filepath.Join("**", "*.js"),
				Confidence: 0.1,
			},
			{
				Type:       PatternDependency,
				Pattern:    "mongoose",
				Path:       "package.json",
				Confidence: 0.1,
			},
		},
	},
}
