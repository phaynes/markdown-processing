FROM rust:1.69 as builder

WORKDIR /usr/src/apa7_to_bibtex
COPY . .
RUN cargo build --release

FROM debian:buster-slim
COPY --from=builder /usr/src/apa7_to_bibtex/target/release/apa7_to_bibtex /usr/local/bin/apa7_to_bibtex
CMD ["apa7_to_bibtex"]
