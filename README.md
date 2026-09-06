# bkdedupe

Every browser lets you export bookmarks to an HTML file but none of them
will tell you which URLs you've saved five times across a decade of synced
profiles, imported backups, and "just in case" exports. bkdedupe reads one
of those export files and prints the URLs that show up more than once, with
the line number and title of each occurrence so you can go clean them up by
hand.

## Usage

Export your bookmarks (Chrome: `chrome://bookmarks` > menu > Export
bookmarks; Firefox: Library > Import and Backup > Export Bookmarks to HTML;
Safari: File > Export Bookmarks), then run:

```sh
bkdedupe bookmarks.html
```

or pipe it in:

```sh
cat bookmarks.html | bkdedupe
```

Example output:

```
https://example.com/pricing (3 times)
  line 42: Example - Pricing
  line 108: Example Pricing (old)
  line 391: Pricing – Example
no duplicate bookmarks found
```

(The second line above only prints when there truly are none — you won't
see both in the same run.)

## Building

```sh
go build -o bkdedupe .
```

## How it works

Netscape-format bookmark files (the format every major browser exports to)
write one `<DT><A HREF="...">Title</A>` per line. bkdedupe reads the file
one line at a time with a buffered reader instead of loading the whole
thing into memory. This matters more than it sounds: browsers often embed
each bookmark's favicon as a base64 data URI right in the `ICON` attribute,
so a bookmark file with a few thousand entries can easily be a few hundred
megabytes even though the actual URLs and titles are tiny.

## Limitations

- Assumes the one-entry-per-line convention that Chrome, Firefox, and
  Safari all follow. Hand-edited or oddly reformatted export files may not
  parse correctly.
- Output order isn't currently sorted (see NEXT ideas below).

## License

MIT, see LICENSE.
