package app.dto;

import jakarta.validation.constraints.NotBlank;

public record JobStatusDocument(
    @NotBlank String repoName, 
    @NotBlank String status
) {}