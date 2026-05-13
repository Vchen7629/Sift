package app.dto;

import java.time.Instant;
import java.util.Map;

import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import com.fasterxml.jackson.databind.ser.std.ToStringSerializer;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;

public record UserRepoDocument (
    @NotBlank String userId,
    @NotBlank String repoName,
    @NotEmpty Map<String, String> dependencies,
    @JsonSerialize(using = ToStringSerializer.class) Instant lastIndexed
) {};