package views

import "tui/internal/types"

var dummyData = []types.Repository{
	{
		GithubId: 0, Name: "react", Status: "indexed", LastIndexed: "1746", TotalDependencies: 42,
		Dependencies: []types.DependencyStatus{
			{Name: "axios",       Version: "v1.7.2",   Status: "healthy"},
			{Name: "lodash",      Version: "v4.17.21", Status: "deprecated"},
			{Name: "moment",      Version: "v2.29.4",  Status: "archived"},
			{Name: "react-query", Version: "v5.28.0",  Status: "healthy"},
			{Name: "classnames",  Version: "v2.3.2",   Status: "healthy"},
			{Name: "prop-types",  Version: "v15.8.1",  Status: "deprecated"},
			{Name: "redux",       Version: "v4.2.1",   Status: "healthy"},
		},
	},
	{
		GithubId: 2, Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalDependencies: 3,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "zod",      Version: "v3.22.4", Status: "healthy"},
			{Name: "immer",    Version: "v10.0.3", Status: "healthy"},
			{Name: "date-fns", Version: "v2.30.0", Status: "deprecated"},
		},
	},
	{
		GithubId: 3, Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalDependencies: 9,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "express",      Version: "v4.19.2", Status: "healthy"},
			{Name: "dotenv",       Version: "v16.4.5", Status: "healthy"},
			{Name: "cors",         Version: "v2.8.5",  Status: "healthy"},
			{Name: "helmet",       Version: "v7.1.0",  Status: "healthy"},
			{Name: "morgan",       Version: "v1.10.0", Status: "deprecated"},
			{Name: "body-parser",  Version: "v1.20.2", Status: "archived"},
			{Name: "jsonwebtoken", Version: "v9.0.2",  Status: "healthy"},
			{Name: "bcrypt",       Version: "v5.1.1",  Status: "healthy"},
			{Name: "multer",       Version: "v1.4.5",  Status: "deprecated"},
		},
	},
	{
		GithubId: 4, Name: "vue", Status: "indexed", LastIndexed: "1746", TotalDependencies: 5,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "pinia",        Version: "v2.1.7",  Status: "healthy"},
			{Name: "vue-router",   Version: "v4.3.0",  Status: "healthy"},
			{Name: "vuelidate",    Version: "v2.0.3",  Status: "deprecated"},
			{Name: "vue-i18n",     Version: "v9.10.2", Status: "healthy"},
			{Name: "vuex",         Version: "v4.1.0",  Status: "archived"},
		},
	},
	{
		GithubId: 5, Name: "svelte", Status: "indexed", LastIndexed: "1745", TotalDependencies: 7,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "svelte-kit",     Version: "v2.5.7",  Status: "healthy"},
			{Name: "svelte-query",   Version: "v5.0.0",  Status: "healthy"},
			{Name: "svelte-forms",   Version: "v1.3.2",  Status: "deprecated"},
			{Name: "svelte-motion",  Version: "v0.12.2", Status: "archived"},
			{Name: "svelte-i18n",    Version: "v3.7.4",  Status: "healthy"},
			{Name: "svelte-persist", Version: "v1.1.0",  Status: "deprecated"},
			{Name: "svelte-meta",    Version: "v0.5.1",  Status: "healthy"},
		},
	},
	{
		GithubId: 6, Name: "angular", Status: "pending", LastIndexed: "0", TotalDependencies: 4,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "rxjs",          Version: "v7.8.1",  Status: "healthy"},
			{Name: "ngrx",          Version: "v17.2.0", Status: "healthy"},
			{Name: "angular-forms", Version: "v16.0.0", Status: "archived"},
			{Name: "zone.js",       Version: "v0.14.4", Status: "deprecated"},
		},
	},
	{
		GithubId: 7, Name: "remix", Status: "indexed", LastIndexed: "1746", TotalDependencies: 6,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "prisma",          Version: "v5.13.0", Status: "healthy"},
			{Name: "zod",             Version: "v3.22.4", Status: "healthy"},
			{Name: "tailwindcss",     Version: "v3.4.3",  Status: "healthy"},
			{Name: "conform",         Version: "v1.1.1",  Status: "healthy"},
			{Name: "remix-auth",      Version: "v3.6.0",  Status: "deprecated"},
			{Name: "react-hot-toast", Version: "v2.4.1",  Status: "archived"},
		},
	},
	{
		GithubId: 8, Name: "astro", Status: "indexed", LastIndexed: "1745", TotalDependencies: 3,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "sharp",      Version: "v0.33.3", Status: "healthy"},
			{Name: "mdx",        Version: "v3.1.0",  Status: "healthy"},
			{Name: "nanostores", Version: "v0.10.3", Status: "deprecated"},
		},
	},
	{
		GithubId: 9, Name: "nuxt", Status: "pending", LastIndexed: "0", TotalDependencies: 8,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "pinia",        Version: "v2.1.7",  Status: "healthy"},
			{Name: "vee-validate",  Version: "v4.12.8", Status: "deprecated"},
			{Name: "nuxt-auth",    Version: "v0.6.7",  Status: "archived"},
			{Name: "nuxt-i18n",    Version: "v8.3.1",  Status: "healthy"},
			{Name: "nuxt-image",   Version: "v1.7.0",  Status: "healthy"},
			{Name: "nuxt-content", Version: "v2.12.1", Status: "healthy"},
			{Name: "vueuse",       Version: "v10.9.0", Status: "healthy"},
			{Name: "ofetch",       Version: "v1.3.4",  Status: "deprecated"},
		},
	},
	{
		GithubId: 10, Name: "solid-js", Status: "indexed", LastIndexed: "1746", TotalDependencies: 5,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "solid-router",  Version: "v0.13.6", Status: "healthy"},
			{Name: "solid-query",   Version: "v5.28.0", Status: "healthy"},
			{Name: "solid-forms",   Version: "v0.4.2",  Status: "deprecated"},
			{Name: "modular-forms", Version: "v0.22.1", Status: "healthy"},
			{Name: "solid-dnd",     Version: "v0.7.4",  Status: "archived"},
		},
	},
	{
		GithubId: 11, Name: "qwik", Status: "indexed", LastIndexed: "1745", TotalDependencies: 7,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "qwik-city",        Version: "v1.5.5",  Status: "healthy"},
			{Name: "builder.io",       Version: "v2.0.1",  Status: "deprecated"},
			{Name: "partytown",        Version: "v0.10.2", Status: "archived"},
			{Name: "mitosis",          Version: "v0.4.6",  Status: "healthy"},
			{Name: "qwik-speak",       Version: "v1.1.0",  Status: "healthy"},
			{Name: "vite-plugin-qwik", Version: "v1.5.5",  Status: "healthy"},
			{Name: "qwik-auth",        Version: "v0.2.0",  Status: "deprecated"},
		},
	},
	{
		GithubId: 12, Name: "preact", Status: "pending", LastIndexed: "0", TotalDependencies: 4,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "preact-router",          Version: "v4.1.2",  Status: "archived"},
			{Name: "preact-signals",         Version: "v1.2.3",  Status: "healthy"},
			{Name: "preact-iso",             Version: "v2.6.2",  Status: "healthy"},
			{Name: "htm",                    Version: "v3.1.1",  Status: "deprecated"},
		},
	},
	{
		GithubId: 13, Name: "lit", Status: "indexed", LastIndexed: "1746", TotalDependencies: 6,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "lit-html",       Version: "v3.1.3",  Status: "healthy"},
			{Name: "lit-element",    Version: "v4.0.4",  Status: "healthy"},
			{Name: "haunted",        Version: "v4.8.1",  Status: "archived"},
			{Name: "shoelace",       Version: "v2.15.0", Status: "healthy"},
			{Name: "lion",           Version: "v0.18.0", Status: "deprecated"},
			{Name: "wired-elements", Version: "v2.0.6",  Status: "archived"},
		},
	},
	{
		GithubId: 14, Name: "ember", Status: "indexed", LastIndexed: "1745", TotalDependencies: 3,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "ember-data",        Version: "v5.3.1",  Status: "deprecated"},
			{Name: "ember-simple-auth", Version: "v6.0.0",  Status: "archived"},
			{Name: "ember-cli-mirage",  Version: "v3.0.4",  Status: "archived"},
		},
	},
	{
		GithubId: 15, Name: "alpine", Status: "pending", LastIndexed: "0", TotalDependencies: 5,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{
			{Name: "alpinejs",             Version: "v3.13.10", Status: "healthy"},
			{Name: "alpine-ajax",          Version: "v0.12.0",  Status: "healthy"},
			{Name: "alpine-magic-helpers", Version: "v0.9.0",   Status: "archived"},
			{Name: "pikaday",              Version: "v1.8.2",   Status: "deprecated"},
			{Name: "sweetalert2",          Version: "v11.10.8", Status: "healthy"},
		},
	},
	{
		GithubId: 16, Name: "no dependencies", Status: "pending", LastIndexed: "0", TotalDependencies: 0,
		Description: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
		Dependencies: []types.DependencyStatus{},
	},
}
