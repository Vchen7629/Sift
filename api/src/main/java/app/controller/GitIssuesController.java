package app.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;

@RestController
@RequestMapping("/issue")
@Validated
public class GitIssuesController {

    private record searchIssueRequest(@NotBlank String query) {}
    
    @PostMapping("/search")
    public ResponseEntity<Void> searchIssue(@RequestBody @Valid searchIssueRequest query) {

        return ResponseEntity.ok().build();
    }
}
