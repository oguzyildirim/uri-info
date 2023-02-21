-- name: SelectURL :one
SELECT * FROM urls
WHERE id = @id LIMIT 1;

-- name: InsertURL :one
INSERT INTO urls (
  HTML_version,
  page_title,
  headings_count,
  links_count,
  inaccessible_links_count,
  have_login_form
)
VALUES (
  @HTMLVersion,
  @pageTitle,
  @headingsCount,
  @linksCount,
  @inaccessibleLinksCount,
  @haveLoginForm
)
RETURNING id;

-- name: DeleteURL :one
DELETE FROM urls
WHERE  id = @id RETURNING id AS res;
