INSERT INTO rule_sets (
  identifier,
  can_delete_any_user_account,
  can_delete_any_thread,
  can_delete_any_post
) VALUES (
  'regular_user_rule_set',
  false,
  false,
  false
);

INSERT INTO site_account_types (
  identifier,
  rule_set_name
) VALUES (
  'regular_user',
  'regular_user_rule_set'
);
