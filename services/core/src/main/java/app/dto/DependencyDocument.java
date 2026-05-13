package app.dto;

import jakarta.validation.constraints.NotBlank;

public record DependencyDocument(
    @NotBlank String name,
    @NotBlank String version,
    @NotBlank String status
) {}
