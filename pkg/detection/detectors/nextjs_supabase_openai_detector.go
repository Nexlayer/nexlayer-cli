// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.

// DEPRECATED: This detector is deprecated and will be removed in a future version.
// The functionality has been replaced by the unified StackDetector in pkg/detection/stack_detector.go,
// which uses pattern-based detection to identify technology stacks.
// Please migrate any direct references to this detector to use StackDetector instead.
// Planned removal: v1.x.0 (next major/minor release)

package detectors

import (
	"context"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// NextjsSupabaseOpenAIDetector detects Next.js + Supabase + OpenAI/Gemini API stack
type NextjsSupabaseOpenAIDetector struct {
	detection.BaseDetector
	nextjsDetector   *NextjsDetector
	supabaseDetector *SupabaseDetector
	openaiDetector   *OpenAIDetector
	geminiDetector   *GeminiDetector
	tailwindDetector *TailwindDetector
	stripeDetector   *StripeDetector
}

// NewNextjsSupabaseOpenAIDetector creates a new detector for Next.js + Supabase + OpenAI/Gemini stack
func NewNextjsSupabaseOpenAIDetector() *NextjsSupabaseOpenAIDetector {
	detector := &NextjsSupabaseOpenAIDetector{
		nextjsDetector:   NewNextjsDetector(),
		supabaseDetector: NewSupabaseDetector(),
		openaiDetector:   NewOpenAIDetector(),
		geminiDetector:   NewGeminiDetector(),
		tailwindDetector: NewTailwindDetector(),
		stripeDetector:   NewStripeDetector(),
	}
	detector.BaseDetector = *detection.NewBaseDetector("Next.js + Supabase + OpenAI/Gemini Detector", 0.8)
	return detector
}

// Detect implementation for NextjsSupabaseOpenAIDetector
func (d *NextjsSupabaseOpenAIDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("NextjsSupabaseOpenAIDetector")

	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Detect individual components
	nextjsInfo, _ := d.nextjsDetector.Detect(ctx, dir)
	supabaseInfo, _ := d.supabaseDetector.Detect(ctx, dir)

	// Check for either OpenAI or Gemini API
	openaiInfo, _ := d.openaiDetector.Detect(ctx, dir)
	geminiInfo, _ := d.geminiDetector.Detect(ctx, dir)

	// Detect optional components
	tailwindInfo, _ := d.tailwindDetector.Detect(ctx, dir)
	stripeInfo, _ := d.stripeDetector.Detect(ctx, dir)

	// Use the higher confidence AI API detection
	aiInfo := openaiInfo
	if geminiInfo.Confidence > openaiInfo.Confidence {
		aiInfo = geminiInfo
	}

	// If we have Next.js + Supabase + AI API (OpenAI or Gemini)
	if nextjsInfo.Type == "nextjs" &&
		supabaseInfo.Type == "supabase" &&
		(openaiInfo.Type != "unknown" || geminiInfo.Type != "unknown") {

		// Calculate combined confidence score
		combinedConf := (nextjsInfo.Confidence + supabaseInfo.Confidence + aiInfo.Confidence) / 3.0

		// Determine the full stack type based on which AI component was detected
		stackType := "nextjs-supabase-openai"
		if geminiInfo.Confidence > openaiInfo.Confidence {
			stackType = "nextjs-supabase-gemini"
		}

		// Set the project info
		projectInfo.Type = stackType
		projectInfo.Confidence = combinedConf
		projectInfo.Language = "JavaScript"
		projectInfo.Framework = "Next.js"
		projectInfo.LLMProvider = aiInfo.LLMProvider

		// Add metadata about detected components
		projectInfo.Metadata["has_tailwind"] = tailwindInfo.Type == "tailwind" && tailwindInfo.Confidence > 0.7
		projectInfo.Metadata["has_stripe"] = stripeInfo.Type == "stripe" && stripeInfo.Confidence > 0.7

		return projectInfo, nil
	}

	// If we only have Next.js + Supabase but no AI component
	if nextjsInfo.Type == "nextjs" && supabaseInfo.Type == "supabase" {
		projectInfo.Type = "nextjs-supabase"
		projectInfo.Confidence = (nextjsInfo.Confidence + supabaseInfo.Confidence) / 2.0
		projectInfo.Language = "JavaScript"
		projectInfo.Framework = "Next.js"

		// Add metadata about detected components
		projectInfo.Metadata["has_tailwind"] = tailwindInfo.Type == "tailwind" && tailwindInfo.Confidence > 0.7
		projectInfo.Metadata["has_stripe"] = stripeInfo.Type == "stripe" && stripeInfo.Confidence > 0.7

		return projectInfo, nil
	}

	return projectInfo, nil
}
