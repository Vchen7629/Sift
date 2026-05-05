package app.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.validation.constraints.NotBlank;

@JsonIgnoreProperties(ignoreUnknown = true)
public record IssueSearchResponse (
    @NotBlank String url,
    @NotBlank String title,
    @NotBlank String body
) {};
