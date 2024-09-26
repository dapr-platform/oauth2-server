dp-cli gen --connstr "postgresql://things:things2024@localhost:5432/thingsdb?sslmode=disable" \
--tables=o_user,o_resource,o_role --model_naming "{{ toUpperCamelCase ( replace . \"o_\" \"\") }}"  \
--file_naming "{{ toLowerCamelCase ( replace . \"o_\" \"\") }}" \
--module oauth2-server

