dp-cli gen --connstr "postgresql://things:things2024@ali4:37432/thingsdb?sslmode=disable" \
--tables=o_user,o_client_info --model_naming "{{ toUpperCamelCase ( replace . \"o_\" \"\") }}"  \
--file_naming "{{ toLowerCamelCase ( replace . \"o_\" \"\") }}" \
--module oauth2-server

