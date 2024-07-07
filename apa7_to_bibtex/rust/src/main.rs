extern crate regex;

use regex::Regex;
use std::env;
use std::fs::File;
use std::io::{BufRead, BufReader, BufWriter, Write};
use std::path::PathBuf;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();
    let input_file = if args.len() > 1 { &args[1] } else { "input.md" };
    let output_file = if args.len() > 2 {
        args[2].clone()
    } else {
        let input_path = PathBuf::from(input_file);
        input_path
            .with_extension("bib")
            .to_string_lossy()
            .into_owned()
    };

    println!("Current directory: {:?}", std::env::current_dir()?);
    println!("Input file: {}", input_file);
    println!("Output file: {}", output_file);

    let file = File::open(input_file).map_err(|e| {
        eprintln!("Error opening input file: {}", e);
        e
    })?;
    let reader = BufReader::new(file);

    let output = File::create(&output_file).map_err(|e| {
        eprintln!("Error creating output file: {}", e);
        e
    })?;
    let mut writer = BufWriter::new(output);

    let author_year_regex = Regex::new(r"^(.*?)\. \((\d{4})\)\. ").unwrap();
    let title_regex = Regex::new(r"\. \*([^\*]+)\*\.").unwrap();
    let publisher_regex = Regex::new(r"\. ([^\.]+)\.$").unwrap();
    let edition_regex = Regex::new(r"\((\d+)(?:th|st|nd|rd) ed\.\)").unwrap();
    let journal_regex = Regex::new(r"\*([^\*]+)\*, (\d+), (\d+[-–]\d+).").unwrap();
    let doi_regex = Regex::new(r"https://doi\.org/([^\s]+)").unwrap();

    for (index, line) in reader.lines().enumerate() {
        let line = line?;
        if line.trim().is_empty() {
            continue;
        }

        let mut entry = String::new();

        if let Some(captures) = author_year_regex.captures(&line) {
            let authors = captures.get(1).unwrap().as_str();
            let year = captures.get(2).unwrap().as_str();
            let first_author = authors
                .split(',')
                .next()
                .unwrap()
                .split_whitespace()
                .last()
                .unwrap()
                .to_lowercase();

            entry.push_str(&format!("@book{{{}_{}},\n", first_author, year));
            entry.push_str(&format!("  author = {{{}}},\n", authors));
            entry.push_str(&format!("  year = {{{}}},\n", year));
        }

        if let Some(captures) = title_regex.captures(&line) {
            let title = captures.get(1).unwrap().as_str();
            entry.push_str(&format!("  title = {{{}}},\n", title));
        }

        if let Some(captures) = edition_regex.captures(&line) {
            let edition = captures.get(1).unwrap().as_str();
            entry.push_str(&format!("  edition = {{{} ed.}},\n", edition));
        }

        if let Some(captures) = publisher_regex.captures(&line) {
            let publisher = captures.get(1).unwrap().as_str();
            entry.push_str(&format!("  publisher = {{{}}},\n", publisher));
        }

        if let Some(captures) = journal_regex.captures(&line) {
            let journal = captures.get(1).unwrap().as_str();
            let volume = captures.get(2).unwrap().as_str();
            let pages = captures.get(3).unwrap().as_str();
            entry.push_str(&format!("  journal = {{{}}},\n", journal));
            entry.push_str(&format!("  volume = {{{}}},\n", volume));
            entry.push_str(&format!("  pages = {{{}}},\n", pages));
        }

        if let Some(captures) = doi_regex.captures(&line) {
            let doi = captures.get(1).unwrap().as_str();
            entry.push_str(&format!("  doi = {{https://doi.org/{}}},\n", doi));
        }

        entry.push_str("}\n");

        writer.write_all(entry.as_bytes())?;
    }

    writer.flush()?;
    Ok(())
}
