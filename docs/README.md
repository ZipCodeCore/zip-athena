# Documentation Index

## Overview

This directory contains the documentation for the Athena project. Documentation is organized to support both users learning from the application and developers contributing to it.

## Main Guides (Located in Root Directory)

These are the primary entry points for different audiences:

- **[Getting Started Guide](../GETTING_STARTED.md)** - Quick setup and first steps
  - Installation instructions
  - Running the application
  - First quiz walkthrough
  - Troubleshooting

- **[Learning Guide](../LEARNING_GUIDE.md)** - Learn software development from this codebase
  - 6 guided learning pathways
  - Code reading techniques
  - Exercises by skill level
  - Study group resources

- **[Architecture Guide](../ARCHITECTURE.md)** - Understand how Athena works
  - System architecture
  - Component details
  - Algorithms (FSRS, quiz generation)
  - Design decisions

- **[Development Guide](../DEVELOPMENT_GUIDE.md)** - Contribute to the project
  - Development workflow
  - Adding features
  - Testing guidelines
  - Code quality standards

- **[README](../README.md)** - Project overview and quick links

## Technical References (`development/`)

Detailed technical documentation for developers:

- **[API_EXAMPLES.md](development/API_EXAMPLES.md)** - Complete API reference
  - All endpoints documented
  - Request/response examples
  - Authentication details
  - Error handling

- **[DATABASE_MIGRATIONS.md](development/DATABASE_MIGRATIONS.md)** - Schema management
  - Migration workflow
  - Creating migrations
  - Rollback procedures
  - Best practices

- **[IMPORT_EXPORT_FORMAT.md](development/IMPORT_EXPORT_FORMAT.md)** - Data interchange
  - JSON format specifications
  - CSV import format
  - Export procedures
  - Validation rules

- **[TEST_SUMMARY.md](development/TEST_SUMMARY.md)** - Testing documentation
  - Test organization
  - Running tests
  - Coverage reports
  - Writing new tests

## Specifications

High-level design and requirements:

- **[SPECIFICATION.md](SPECIFICATION.md)** - Complete system specification
  - Feature requirements
  - Technical specifications
  - User stories
  - System constraints

- **[CURRICULUM.md](CURRICULUM.md)** - Learning content organization
  - Topics covered
  - Skill progression
  - Resource mapping
  - Learning paths

## Development History (`archive/`)

Historical documentation from the development process:

- Phase implementation documents (PHASE1-5)
- Feature implementation notes
- Completion summaries
- Development process artifacts

These documents are preserved for reference but are not the primary documentation for understanding the current system. See the main guides above instead.

## Documentation Philosophy

### For Learners

We focus documentation on:
1. **Understanding** - What the code does and why
2. **Learning** - Using the codebase as an educational resource
3. **Exploration** - Multiple pathways to explore topics of interest

### For Developers

We emphasize:
1. **Getting Started** - Quick setup and contribution
2. **Best Practices** - Code quality and standards
3. **Technical Details** - APIs, schemas, and internals

## Finding What You Need

**I want to...**

- **Run the app** → [Getting Started Guide](../GETTING_STARTED.md)
- **Learn from the code** → [Learning Guide](../LEARNING_GUIDE.md)
- **Understand architecture** → [Architecture Guide](../ARCHITECTURE.md)
- **Contribute code** → [Development Guide](../DEVELOPMENT_GUIDE.md)
- **Use the API** → [API Examples](development/API_EXAMPLES.md)
- **Modify the schema** → [Database Migrations](development/DATABASE_MIGRATIONS.md)
- **Import questions** → [Import/Export Format](development/IMPORT_EXPORT_FORMAT.md)
- **Write tests** → [Test Summary](development/TEST_SUMMARY.md)

## Contributing to Documentation

Documentation improvements are always welcome! When contributing:

1. **Keep user perspective** - Who is reading this and what do they need?
2. **Use examples** - Show, don't just tell
3. **Update related docs** - Keep cross-references current
4. **Test instructions** - Verify commands actually work
5. **Be clear and concise** - Respect the reader's time

### Documentation Standards

- Use Markdown for all documentation
- Include code examples with syntax highlighting
- Keep line length reasonable (80-100 characters)
- Use relative links for internal references
- Include table of contents for longer documents

## Feedback

Found unclear documentation? Please:
1. Open an issue describing what's confusing
2. Suggest improvements
3. Submit a PR with updates

Good documentation makes the project accessible to everyone!
