package app.exception;

import java.io.IOException;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

import jakarta.validation.ConstraintViolationException;

@ControllerAdvice
public class GeneralExceptionHandler {
    @ExceptionHandler(IOException.class)
    public ResponseEntity<String> IOException(IOException e) {
        return ResponseEntity.status(500).body(e.getMessage());
    }

    @ExceptionHandler(ConstraintViolationException.class)
    public ResponseEntity<String> constraintViolationException(IOException e) {
        return ResponseEntity.status(500).body("internal server error");
    }
}
