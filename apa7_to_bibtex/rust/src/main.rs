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
    let title_regex = Regex::new(r"\. (.*?)\.").unwrap();
    let journal_regex = Regex::new(r"\. \*(.*?)\*").unwrap();
    let volume_issue_regex = Regex::new(r", (\d+)\((\d+)\)").unwrap();
    let pages_regex = Regex::new(r", (\d+[-–]\d+)").unwrap();
    let doi_regex = Regex::new(r"https://doi\.org/(.*)").unwrap();

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
                .trim()
                .split(' ')
                .last()
                .unwrap();
            let key = format!("{}_{}", first_author.to_lowercase(), year);

            entry.push_str(&format!("@article{{{},\n", key));
            entry.push_str(&format!("  author = {{{}}},\n", authors));
            entry.push_str(&format!("  year = {{{}}},\n", year));
        } else {
            eprintln!(
                "Warning: Line {} does not match the expected format for author and year.",
                index + 1
            );
            continue;
        }

        if let Some(captures) = title_regex.captures(&line) {
            let title = captures.get(1).unwrap().as_str();
            entry.push_str(&format!("  title = {{{}}},\n", title));
        } else {
            eprintln!("Warning: Line {} is missing a title.", index + 1);
        }

        if let Some(captures) = journal_regex.captures(&line) {
            let journal = captures.get(1).unwrap().as_str();
            entry.push_str(&format!("  journal = {{{}}},\n", journal));
        } else {
            eprintln!("Warning: Line {} is missing a journal.", index + 1);
        }

        if let Some(captures) = volume_issue_regex.captures(&line) {
            let volume = captures.get(1).unwrap().as_str();
            let issue = captures.get(2).unwrap().as_str();
            entry.push_str(&format!("  volume = {{{}}},\n", volume));
            entry.push_str(&format!("  number = {{{}}},\n", issue));
        } else {
            eprintln!(
                "Warning: Line {} is missing volume and issue information.",
                index + 1
            );
        }

        if let Some(captures) = pages_regex.captures(&line) {
            let pages = captures.get(1).unwrap().as_str();
            entry.push_str(&format!("  pages = {{{}}},\n", pages));
        } else {
            eprintln!("Warning: Line {} is missing page information.", index + 1);
        }

        if let Some(captures) = doi_regex.captures(&line) {
            let doi = captures.get(1).unwrap().as_str();
            entry.push_str(&format!("  doi = {{{}}}\n", doi));
        } else {
            eprintln!("Warning: Line {} is missing a DOI.", index + 1);
        }

        entry.push_str("}\n\n");
        writer.write_all(entry.as_bytes())?;
    }

    writer.flush()?;
    println!("Conversion complete. BibTeX file created: {}", output_file);
    Ok(())
}
