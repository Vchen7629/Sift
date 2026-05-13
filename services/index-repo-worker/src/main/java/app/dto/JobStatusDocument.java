package app.dto;

import jakarta.validation.constraints.NotBlank;

public record JobStatusDocument(
    @NotBlank String userId,
    @NotBlank String repoName, 
    @NotBlank String status
) {}