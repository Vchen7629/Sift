package app.dto;

import jakarta.validation.constraints.NotBlank;

public record GithubChangeLogResponse (
    @NotBlank String dependencyName,
    @NotBlank String version, 
    @NotBlank String changes,
    @NotBlank String url
) {}