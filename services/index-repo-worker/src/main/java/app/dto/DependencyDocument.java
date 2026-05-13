package app.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.validation.constraints.NotBlank;

@JsonIgnoreProperties(ignoreUnknown = true)
public record DependencyDocument(
    @NotBlank String dependencyName, 
    @NotBlank String version
) {}
