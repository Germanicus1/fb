# Flow Boards CLI

`fb` is a terminal client for Flow Boards: it shows the tickets assigned to you and lets you comment on them without leaving the shell. Most of its vocabulary is inherited from Flow Boards; the checkout workflow is the CLI's own invention, layered on top.

## Flow Boards concepts

**Ticket**:
A single unit of work in Flow Boards, assigned to one or more users. Identified by an opaque alphanumeric ID.
_Avoid_: issue, card, task, item

**Bin**:
A named container on a board that holds tickets. A ticket lives in exactly one bin.
_Avoid_: column, list, lane, stage, bucket

**Board**:
A named collection of bins.
_Avoid_: project, workspace, swimlane

**Status**:
A ticket's bin, rendered for display — its bin name, falling back to bin ID, then `Unknown`. Never a field of its own on a ticket.
_Avoid_: state, stage

**Comment**:
A body of text posted against a ticket. Write-only from this CLI — `fb` posts comments but never reads them back.
_Avoid_: note, message, update

**REST prefix**:
The organization-specific API base URL, discovered from the org ID before any other request.
_Avoid_: endpoint, base URL, API host

## Checkout workflow

**Checkout**:
The one ticket currently pinned as the target of quick comments, persisted in `~/.fb/checkout.json` so it survives across shell sessions. At most one exists at a time.
_Avoid_: pin, lock, selection, active ticket, current ticket

**Checkout indicator**:
The `← CHECKED OUT` suffix appended to the checked-out ticket's line in a listing.
_Avoid_: marker, badge, highlight, flag

**Quick comment**:
A comment posted straight to the checkout with no ticket selection step (`fb -c "…"`), as opposed to the interactive `fb --comment`.
_Avoid_: fast comment, instant comment, one-liner

**Bin context**:
The bin last used for a checkout, remembered in `~/.fb/bin_context.json` so a bare `fb checkout` knows where to look.
_Avoid_: last bin, default bin, sticky bin, remembered bin

## Filtering and output

**Bin filter**:
The `--bin` argument, accepted as either a bin ID or a bin name. Resolving it means deciding which of the two it is (IDs are alphanumeric; names carry spaces or punctuation) and turning a name into an ID.
_Avoid_: bin query, bin selector, bin argument

**Server-side filtering**:
Narrowing tickets by passing bin/board IDs to the API, so the response only carries what's wanted. Preferred over **client-side filtering**, which trims an already-fetched slice and exists mainly for names the API can't resolve.
_Avoid_: remote filtering, API filtering, pre-filtering

**Minimal output**:
The default ticket listing — ID and name only, one line per ticket.
_Avoid_: short output, quiet mode, compact mode

**Verbose output**:
The `--verbose` / `-v` / `--debug` listing — full ticket detail (status, dates, wrapped description) plus API and total timings.
_Avoid_: debug mode, detailed mode, full mode
