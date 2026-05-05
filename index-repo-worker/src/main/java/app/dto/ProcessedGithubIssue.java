package app.dto;

import java.util.List;

public record ProcessedGithubIssue(
    String dependencyName, 
    String version,
    String title, 
    String body,
    String url, 
    List<String> labelList,
    String createdOn
) {}