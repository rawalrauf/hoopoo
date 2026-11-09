.PHONY: all clean install test bench

# Build all Go binaries
all:
	@echo "Building Go binaries..."
	@mkdir -p bin
	@for cmd in cmd/hoopoo-*; do \
		if [ -d "$$cmd" ]; then \
			name=$$(basename $$cmd); \
			echo "  Building $$name..."; \
			go build -ldflags="-s -w" -o bin/$$name ./$$cmd; \
		fi \
	done

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f bin/hoopoo-* 2>/dev/null || true
	@find bin -type f -executable -exec file {} \; | grep ELF | cut -d: -f1 | xargs rm -f

# Install to /usr/bin (requires sudo)
install: all
	@echo "Installing to /usr/bin..."
	@for bin in bin/hoopoo-*; do \
		if [ -f "$$bin" ] && file "$$bin" | grep -q ELF; then \
			sudo install -Dm755 "$$bin" /usr/bin/$$(basename $$bin); \
		fi \
	done

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@echo "Benchmark targets removed (stale battery commands cleaned up)"

# Test a single command
test:
	@echo "Testing Go binaries..."
	@for bin in bin/hoopoo-*; do \
		if [ -f "$$bin" ] && file "$$bin" | grep -q ELF; then \
			echo "Testing $$(basename $$bin)..."; \
			$$bin || true; \
		fi \
	done
