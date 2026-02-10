# Phase 2 and 3 Learning Resources Fix

## Problem

Phase 2 and Phase 3 learning resources were showing in the Learning Center outline but returning 404 errors when clicked. This was caused by two issues:

1. **Incorrect file paths in database**: The `file_path` column had incorrect values (e.g., `java/8-classes-and-objects.md` instead of `java/8-oop-classes-and-objects.md`)
2. **Resources marked as inactive**: All Phase 2 and 3 resources had `is_active = 0`

## Files Fixed

### Database Seed File
- **File**: `database/seed_learning_resources.sql`
- **Changes**: 
  - Corrected file paths for all Phase 2 and 3 resources
  - Set `is_active = 1` for all Phase 2 and 3 resources

### Java Phase 2 Resources (OOP)
- `java/8-oop-classes-and-objects.md` (was: `java/8-classes-and-objects.md`)
- `java/9-oop-inheritance.md` (was: `java/9-inheritance.md`)
- `java/10-oop-polymorphism.md` (was: `java/10-polymorphism.md`)
- `java/11-interfaces.md` ✓ (already correct)

### Java Phase 3 Resources (Advanced Data Handling)
- `java/12-packages.md` ✓ (already correct)
- `java/13-exception-handling.md` ✓ (already correct)
- `java/14-file-io.md` ✓ (already correct)
- `java/15-collections-framework.md` (was: `java/15-collections.md`)

### Python Phase 2 Resources (OOP)
- `python/14-oop-classes-and-objects.md` (was: `python/14-classes-and-objects.md`)
- `python/15-oop-inheritance.md` (was: `python/15-inheritance.md`)
- `python/16-oop-polymorphism.md` (was: `python/16-polymorphism.md`)

### Python Phase 3 Resources (Advanced Data Handling)
- `python/11-modules-and-packages.md` (was: `python/11-modules-packages.md`)
- `python/12-file-io.md` ✓ (already correct)
- `python/13-exception-handling.md` ✓ (already correct)

## How to Apply the Fix on Production

### Option 1: Run the Fix Script (Recommended)

```bash
# On your production server
./fix_learning_resources.sh athena.db
```

This script will:
1. Show the current state of Phase 2 and 3 resources
2. Apply the migration
3. Show the updated state
4. Verify that all files exist

### Option 2: Apply Migration Manually

```bash
# On your production server
sqlite3 athena.db < database/migrations/005_fix_phase2_phase3_learning_resources.sql
```

### Option 3: Re-run Complete Setup

If you want to completely reset the learning resources table:

```bash
# On your production server
./setup_learning.sh
```

**Warning**: This will delete and recreate the entire learning_resources table.

## Verification

After applying the fix, verify it works:

1. Open the Learning Center in your browser
2. Click on any Phase 2 or Phase 3 resource
3. The content should now load properly without 404 errors

You can also verify from the command line:

```bash
# Check database records
sqlite3 athena.db "SELECT id, language, phase, section_number, section_title, file_path, is_active FROM learning_resources WHERE phase IN (2, 3) ORDER BY phase, language, section_number;"

# Verify files exist
ls -la learning/java/8-oop-classes-and-objects.md
ls -la learning/java/15-collections-framework.md
ls -la learning/python/14-oop-classes-and-objects.md
ls -la learning/python/11-modules-and-packages.md
```

## Files Created/Modified

1. `database/seed_learning_resources.sql` - Fixed file paths and activated resources
2. `database/migrations/005_fix_phase2_phase3_learning_resources.sql` - Migration to update existing databases
3. `fix_learning_resources.sh` - Helper script to apply the migration
4. `LEARNING_RESOURCES_FIX.md` - This documentation

## Deployment

To deploy this fix:

```bash
# Commit the changes
git add database/seed_learning_resources.sql
git add database/migrations/005_fix_phase2_phase3_learning_resources.sql
git add fix_learning_resources.sh
git add LEARNING_RESOURCES_FIX.md
git commit -m "Fix Phase 2 and 3 learning resource file paths and activate them"
git push

# On production server
git pull
./fix_learning_resources.sh athena.db
```

## Testing Locally

To test this fix locally first:

```bash
# Run the fix script
./fix_learning_resources.sh athena.db

# Start the server
go run main.go

# Open http://localhost:9206/learning.html
# Navigate to any Phase 2 or 3 resource and verify it loads
```
