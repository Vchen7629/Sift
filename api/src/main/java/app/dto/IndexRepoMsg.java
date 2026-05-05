package app.dto;

import jakarta.validation.constraints.NotBlank;

public record IndexRepoMsg(@NotBlank String repoName) {};
