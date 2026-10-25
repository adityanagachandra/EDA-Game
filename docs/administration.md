# Administration, roles, permissions and the audit trail

Mathua has one authentication system and one role boundary. This document is the design
rationale for that boundary: why the role lives in the database and not in the token, why the
permission model is data, why staff arrive by invitation and why the last owner cannot be
removed. It is analytical prose, not a procedure.

For the system overview, see [`architecture.md`](architecture.md).

## Administration, roles, not a second login

One authentication system serves everyone. A learner signs in at `/login`. An administrator
signs in at the same place with the same form. There is no admin login page and no admin
password. A role on the account separates the two.

### The roles

`students.role` holds one of four values. The table lists them from most to least privilege.

| Role | What it can do |
|---|---|
| `owner` | Everything. The only role that can grant another owner. |
| `admin` | Run administration and moderation. Cannot grant owner. |
| `moderator` | Work the moderation queue and inspect content. |
| `student` | Learn. No administrative access. |

The column is `NOT NULL DEFAULT 'student'`. Every account created before roles migrated to a
learner. A database from before roles holds no evidence about who administered it, because
nothing in that database permitted administration.

V1 held only `student` and `admin`. The migration moves each V1 `admin` to `owner`. V1 `admin`
was full control. The move keeps that control. Leaving the value as `admin` takes away the power to
manage staff.

The role never comes from XP, from mastery, from level or from an email domain. A learner can
change each of those. A request body cannot set the role.

### Where the role is read

`requirePermission` (`internal/server/adminrole.go`) is the boundary. It authenticates the
request. It identifies the account. Then it reads the current role from the database. It does
not trust a role carried in the token.

That is the decision the rest of the design follows from. A role embedded in a JWT is a claim
made at sign-in. The claim stays valid until the token expires. Expiry can be a day. An
administrator demoted at 09:00 keeps access until midnight under that scheme. That administrator
can use the access to promote themselves back.

One indexed read per request removes the window. A demotion then takes effect on the next call,
from any device, with no re-login.
`TestLastAdminDemotionTakesEffectImmediately` in `internal/server/adminrole_test.go` asserts
this.

### The permission model

`internal/admin/admin.go` holds the model. Each role maps to a set of permissions. A route names
one permission. The boundary checks the role against that permission.

The model is data, not scattered checks. A permission matrix spread across handlers is a matrix
nobody can review. The failure it produces is a route that forgot a check.

The matrix is small on purpose. `owner` holds every permission. `admin` holds all but
`roles.grant_owner`. `moderator` holds the moderation and content permissions. A learner holds
nothing.

`CanAssign` decides one role change. It applies three rules. An actor cannot grant a role above
their own. An actor cannot act on a peer or a superior, unless the actor is an owner. Granting
`owner` needs `roles.grant_owner`. `TestCanAssign` in `internal/admin/admin_test.go` states the
rules.

### Moderation under the same boundary

Reported content is moderated at `/api/admin/reports`. The role boundary guards it. A moderator
resolves a decision and the same transaction writes the audit row.

V1 gated content triage with a shared `ADMIN_PASSWORD`. V2 removed that mechanism. A shared
password cannot say which person used it. One moderation path replaces the two.

A report is never deleted. A decision sets the status, the reason, the actor and the time. The
row keeps the complaint, so the next moderator can read what was decided before.

A contact message can carry one screenshot. The image is stored with the report and served only
through the same moderation permission, so it is not a public file and it disappears with the
account that filed it.

### Staff access by invitation

An owner or admin adds a contributor from `/admin/contributors`. The server creates an
invitation. The invitation carries a token, an email and a role.

The server stores only a SHA-256 hash of the token. A database read therefore yields no usable
credential. The raw token returns once, in the response and nowhere else.

The invitation is single-use, time-limited and bound to the invited address. Acceptance burns
the token. Acceptance checks the signed-in account against the invited email. A revoked,
expired, reused or mismatched invitation grants nothing. Each failure reads the same to the
caller, so a prober learns nothing about which near-miss they reached.

`ConsumeAdminInvitation` burns the token, grants the role and writes the audit row in one
transaction. No state exists in which someone holds a staff role with no record of how they got
it.

### The last owner

The system never reaches zero owners through an administrative action. The check sits inside the
UPDATE statement. It does not sit before the UPDATE.

```sql
UPDATE students SET role = ?
 WHERE id = ? AND role = 'owner'
   AND (SELECT COUNT(*) FROM students WHERE role = 'owner') > 1
```

A count read before a write loses a race. Two owners who demote each other at the same moment
both read two. Both write. The system reaches zero through the feature meant to stop it. One
statement means the second attempt matches no rows.

The invariant is on the owner, not the admin. You can remove an admin or a moderator freely. An
owner cannot be, because nothing but an owner can create another owner.

### Bootstrap

The operator creates the first owner on the machine that holds the database:

```bash
mathua admin promote <email>
```

The command requires an account that already exists. It refuses to create one. Shell access to
the machine already grants control of the deployment, so the machine is the correct place to
decide who administers.

The command grants `owner`, not `admin`. The first promotion must be able to create the rest.
The command calls the same `SetRoleAudited` as the HTTP handler, with `cli` as the actor. The
command and the web page therefore share one code path that writes the audit row.

An HTTP endpoint fails this test. Anyone who can reach the server could then ask for
administrator access. This repository stores no administrator identity. It stores the mechanism.

### The audit log

Every role change and every moderation decision writes one row to `admin_audit`. The insert and
the change share one transaction. A failed insert rolls the change back. A privileged change
therefore cannot exist without a record.
`TestFailedAuditInsertRollsBackPromotion` in `internal/server` proves this.

The log has no foreign key on `actor_id`. A trail deleted together with the account it records
is not a trail. Neither store offers an update path. Neither store offers a delete path. No HTTP
route offers either. Each event names the actor, the action, the entity, the time and the state
either side of the change.

`ListAdminAuditFor` reads the events about one entity. A role history and a report history are
queries over the trail, not columns anyone maintains by hand.

The before and after snapshots come from `storage.AdminUser`. That type has no field for a
password hash, a share token or learner progress. A credential cannot reach the log through it.

### What the frontend is not

`AdminNav` shows a section only when the caller holds the permission for it. The entry is a
convenience. A hidden link is not authorization.

Every administrative endpoint rejects a caller without the permission on its own. A learner who
types `/admin/users` reaches the page, and then the API refuses the request with a 403.

### Confirmation before a role change

`/admin/users` offers a role picker on each row. The picker opens `ConfirmDialog`. The dialog
names the account and the new role.

The dialog puts focus on `Cancel`. A stray Enter therefore does nothing. Escape, the backdrop
and `Cancel` all close the dialog and send no request.

The page shows two errors apart. One is the list failing to load. The other is the server
refusing a change. They are different sentences.

### Seeing the leaguesThe public league read answers one question: where a learner sits. It cannot answer "what does
this week's pyramid look like", so an administrator had no screen on which a cohort existed. A
report saying "my league is broken" could not be located.

`GET /api/admin/leagues` returns every league of a week, in tier order, with its size and its
mover count. `?tier=bronze&league=3` adds the standings for that one league. The pyramid read is
a single grouped query, and a roster is read only when a league is opened.

The route is guarded by `users.read`. A league row is an account plus a week, so the grant that
lets an administrator look an account up is the grant that lets them enumerate the same people
by cohort.

A week key is a Monday. The handler validates it, because the value is a join key into
`league_members`. A read of the current week settles that week first, exactly as a learner's own
read does. A read of a past week writes nothing.

### The deployed surface

The static export writes both `/admin.html` and a directory named `/admin`. `nextStaticFS` in
`cmd/mathua/main.go` probes the requested path and serves a file when one exists. It must also
require that the path is not a directory. Without that test the server answers `/admin` with a
redirect to `/admin/`, and `/admin/` returns a generated directory listing. Every route answers
`200` in that state, so no page-level test can see the fault.
`cmd/mathua/staticfs_test.go` asserts that no response body holds a directory listing.

### What the machine can carry, measured

The app once ran on `shared-cpu-1x` with 256 MB. On 2026-10-23 that machine stopped answering
for minutes at a time: `/api/health`, which touches no database and writes a constant JSON body,
was logged at 17 to 32 seconds, and the Fly proxy answered 94 requests in 45 seconds with
`could not find a good candidate within 40 attempts at load balancing`, because the 2 second
health check could not get a reply. Requests that did reach the app took 1m45s and the grading
deadline turned them into `503`.

The cause was memory, not processor. The machine reported `MemTotal` 212 MB and, during the
burst, 8 MB available. The Go process reached 101 MB resident, the SymPy child holds 57 MB, and
`/data/mathua.db` is 368 MB. With no memory left for page cache the kernel evicted the process's
own pages to a volume that measures 20 MB/s sequential read, so the runtime spent its time
faulting pages back in. `/proc/pressure/io` reported `full avg10=29.77`, meaning every runnable
task was stalled on I/O for 30% of wall time, while `/proc/pressure/cpu` stayed near 4%. The
processor was idle and the disk was the queue.

The same software is fast when memory is not the constraint. On a 1 CPU / 256 MB container with a
fresh database and rate limits relaxed (`MATHUA_LOADTEST=1`), driving the real diagnostic path,
the machine sustained about 235 answers per second with a median of 1 ms and a 95th percentile of
236 ms at 16 concurrent learners and 459 ms at 32, with no errors. A numeric-only mix reached
about 600 answers per second. So the ceiling at classroom scale is memory and I/O, not the
handler.

Three changes follow from that, and each one is in the repository rather than in a setting:

- `fly.toml` requests 1 GB, which leaves room for page cache beside the heap and the CAS child.
- `GOMEMLIMIT=700MiB` keeps the Go heap from taking the whole machine. Without a limit the
  collector's default target is twice the live heap, which is free to consume the box.
- SQLite runs with a 32 MB page cache and `synchronous=NORMAL`. A commit no longer fsyncs a
  20 MB/s volume, and the cost is the last few transactions if the host loses power, not
  corruption. The client already retries an unsaved answer.

Two related limits are worth knowing. A branch deploy restarts the only machine, and a volume can
be attached to one machine at a time, so a deploy is a short outage rather than a rolling
replacement; the deploy of 2026-10-23T19:08 was visible as that error burst. And symbolic grading
is serialized through one Python pipe, so the wait for that pipe is now bounded at five seconds
and the caller is told to retry rather than queueing behind a slow expression.
