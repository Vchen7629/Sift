package app.exception;

import ai.djl.MalformedModelException;
import ai.djl.repository.zoo.ModelNotFoundException;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice
public class ModelExceptionHandler {
    @ExceptionHandler(ModelNotFoundException.class)
    public ResponseEntity<String> handleModelNotFound(ModelNotFoundException e) {
        return ResponseEntity.status(500).body("embedding model not found/loaded");
    }

    @ExceptionHandler(MalformedModelException.class)
    public ResponseEntity<String> handleMalformedModel(MalformedModelException e) {
        return ResponseEntity.status(500).body("embedding model is malformed or corrupted");
    }
}
