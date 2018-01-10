/* $OpenBSD: tls_config.c,v 1.36 2017/01/31 16:18:57 beck Exp $ */
/*
 * Copyright (c) 2014 Joel Sing <jsing@openbsd.org>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

#include <sys/stat.h>

#include <ctype.h>
#include <errno.h>
#include <fcntl.h>
#include <stdlib.h>
#include <unistd.h>

#include <tls.h>
#include "tls_internal.h"

static int
set_string(const char **dest, const char *src)
{
	free((char *)*dest);
	*dest = NULL;
	if (src != NULL)
		if ((*dest = strdup(src)) == NULL)
			return -1;
	return 0;
}

static void *
memdup(const void *in, size_t len)
{
	void *out;

	if ((out = malloc(len)) == NULL)
		return NULL;
	memcpy(out, in, len);
	return out;
}

static int
set_mem(char **dest, size_t *destlen, const void *src, size_t srclen)
{
	free(*dest);
	*dest = NULL;
	*destlen = 0;
	if (src != NULL)
		if ((*dest = memdup(src, srclen)) == NULL)
			return -1;
	*destlen = srclen;
	return 0;
}

static struct tls_keypair *
tls_keypair_new(void)
{
	return calloc(1, sizeof(struct tls_keypair));
}

static int
tls_keypair_set_cert_file(struct tls_keypair *keypair, struct tls_error *error,
    const char *cert_file)
{
	return tls_config_load_file(error, "certificate", cert_file,
	    &keypair->cert_mem, &keypair->cert_len);
}

static int
tls_keypair_set_cert_mem(struct tls_keypair *keypair, const uint8_t *cert,
    size_t len)
{
	return set_mem(&keypair->cert_mem, &keypair->cert_len, cert, len);
}

static int
tls_keypair_set_key_file(struct tls_keypair *keypair, struct tls_error *error,
    const char *key_file)
{
	if (keypair->key_mem != NULL)
		explicit_bzero(keypair->key_mem, keypair->key_len);
	return tls_config_load_file(error, "key", key_file,
	    &keypair->key_mem, &keypair->key_len);
}

static int
tls_keypair_set_key_mem(struct tls_keypair *keypair, const uint8_t *key,
    size_t len)
{
	if (keypair->key_mem != NULL)
		explicit_bzero(keypair->key_mem, keypair->key_len);
	return set_mem(&keypair->key_mem, &keypair->key_len, key, len);
}

static int
tls_keypair_set_ocsp_staple_file(struct tls_keypair *keypair,
    struct tls_error *error, const char *ocsp_file)
{
	return tls_config_load_file(error, "ocsp", ocsp_file,
	    &keypair->ocsp_staple, &keypair->ocsp_staple_len);
}

static int
tls_keypair_set_ocsp_staple_mem(struct tls_keypair *keypair,
    const uint8_t *staple, size_t len)
{
	return set_mem(&keypair->ocsp_staple, &keypair->ocsp_staple_len, staple,
	    len);
}

static void
tls_keypair_clear(struct tls_keypair *keypair)
{
	tls_keypair_set_cert_mem(keypair, NULL, 0);
	tls_keypair_set_key_mem(keypair, NULL, 0);
}

static void
tls_keypair_free(struct tls_keypair *keypair)
{
	if (keypair == NULL)
		return;

	tls_keypair_clear(keypair);

	free(keypair->cert_mem);
	free(keypair->key_mem);
	free(keypair->ocsp_staple);

	free(keypair);
}

int
tls_config_load_file(struct tls_error *error, const char *filetype,
    const char *filename, char **buf, size_t *len)
{
	struct stat st;
	int fd = -1;
	ssize_t n;

	free(*buf);
	*buf = NULL;
	*len = 0;

	if ((fd = open(filename, O_RDONLY)) == -1) {
		tls_error_set(error, "failed to open %s file '%s'",
		    filetype, filename);
		goto fail;
	}
	if (fstat(fd, &st) != 0) {
		tls_error_set(error, "failed to stat %s file '%s'",
		    filetype, filename);
		goto fail;
	}
	if (st.st_size < 0)
		goto fail;
	*len = (size_t)st.st_size;
	if ((*buf = malloc(*len)) == NULL) {
		tls_error_set(error, "failed to allocate buffer for "
		    "%s file", filetype);
		goto fail;
	}
	n = read(fd, *buf, *len);
	if (n < 0 || (size_t)n != *len) {
		tls_error_set(error, "failed to read %s file '%s'",
		    filetype, filename);
		goto fail;
	}
	close(fd);
	return 0;

 fail:
	if (fd != -1)
		close(fd);
	if (*buf != NULL)
		explicit_bzero(*buf, *len);
	free(*buf);
	*buf = NULL;
	*len = 0;

	return -1;
}

struct tls_config *
tls_config_new(void)
{
	struct tls_config *config;
	unsigned char sid[TLS_MAX_SESSION_ID_LENGTH];

	if ((config = calloc(1, sizeof(*config))) == NULL)
		return (NULL);

	if ((config->keypair = tls_keypair_new()) == NULL)
		goto err;

	/*
	 * Default configuration.
	 */
	if (tls_config_set_dheparams(config, "none") != 0)
		goto err;
	if (tls_config_set_ecdhecurve(config, "auto") != 0)
		goto err;
	if (tls_config_set_ciphers(config, "secure") != 0)
		goto err;

	if (tls_config_set_protocols(config, TLS_PROTOCOLS_DEFAULT) != 0)
		goto err;
	if (tls_config_set_verify_depth(config, 6) != 0)
		goto err;

	/*
	 * Set session ID context to a random value.  For the simple case
	 * of a single process server this is good enough. For multiprocess
	 * servers the session ID needs to be set by the caller.
	 */
	arc4random_buf(sid, sizeof(sid));
	if (tls_config_set_session_id(config, sid, sizeof(sid)) != 0)
		goto err;
	config->ticket_keyrev = arc4random();
	config->ticket_autorekey = 1;

	tls_config_prefer_ciphers_server(config);

	tls_config_verify(config);

	return (config);

 err:
	tls_config_free(config);
	return (NULL);
}

void
tls_config_free(struct tls_config *config)
{
	struct tls_keypair *kp, *nkp;

	if (config == NULL)
		return;

	for (kp = config->keypair; kp != NULL; kp = nkp) {
		nkp = kp->next;
		tls_keypair_free(kp);
	}

	free(config->error.msg);

	free(config->alpn);
	free((char *)config->ca_mem);
	free((char *)config->ca_path);
	free((char *)config->ciphers);

	free(config);
}

static void
tls_config_keypair_add(struct tls_config *config, struct tls_keypair *keypair)
{
	struct tls_keypair *kp;

	kp = config->keypair;
	while (kp->next != NULL)
		kp = kp->next;

	kp->next = keypair;
}

const char *
tls_config_error(struct tls_config *config)
{
	return config->error.msg;
}

void
tls_config_clear_keys(struct tls_config *config)
{
	struct tls_keypair *kp;

	for (kp = config->keypair; kp != NULL; kp = kp->next)
		tls_keypair_clear(kp);

	tls_config_set_ca_mem(config, NULL, 0);
}

int
tls_config_parse_protocols(uint32_t *protocols, const char *protostr)
{
	uint32_t proto, protos = 0;
	char *s, *p, *q;
	int negate;

	if ((s = strdup(protostr)) == NULL)
		return (-1);

	q = s;
	while ((p = strsep(&q, ",:")) != NULL) {
		while (*p == ' ' || *p == '\t')
			p++;

		negate = 0;
		if (*p == '!') {
			negate = 1;
			p++;
		}

		if (negate && protos == 0)
			protos = TLS_PROTOCOLS_ALL;

		proto = 0;
		if (strcasecmp(p, "all") == 0 ||
		    strcasecmp(p, "legacy") == 0)
			proto = TLS_PROTOCOLS_ALL;
		else if (strcasecmp(p, "default") == 0 ||
		    strcasecmp(p, "secure") == 0)
			proto = TLS_PROTOCOLS_DEFAULT;
		if (strcasecmp(p, "tlsv1") == 0)
			proto = TLS_PROTOCOL_TLSv1;
		else if (strcasecmp(p, "tlsv1.0") == 0)
			proto = TLS_PROTOCOL_TLSv1_0;
		else if (strcasecmp(p, "tlsv1.1") == 0)
			proto = TLS_PROTOCOL_TLSv1_1;
		else if (strcasecmp(p, "tlsv1.2") == 0)
			proto = TLS_PROTOCOL_TLSv1_2;

		if (proto == 0) {
			free(s);
			return (-1);
		}

		if (negate)
			protos &= ~proto;
		else
			protos |= proto;
	}

	*protocols = protos;

	free(s);

	return (0);
}

static int
tls_config_parse_alpn(struct tls_config *config, const char *alpn,
    char **alpn_data, size_t *alpn_len)
{
	size_t buf_len, i, len;
	char *buf = NULL;
	char *s = NULL;
	char *p, *q;

	free(*alpn_data);
	*alpn_data = NULL;
	*alpn_len = 0;

	if ((buf_len = strlen(alpn) + 1) > 65535) {
		tls_config_set_errorx(config, "alpn too large");
		goto err;
	}

	if ((buf = malloc(buf_len)) == NULL) {
		tls_config_set_errorx(config, "out of memory");
		goto err;
	}

	if ((s = strdup(alpn)) == NULL) {
		tls_config_set_errorx(config, "out of memory");
		goto err;
	}

	i = 0;
	q = s;
	while ((p = strsep(&q, ",")) != NULL) {
		if ((len = strlen(p)) == 0) {
			tls_config_set_errorx(config,
			    "alpn protocol with zero length");
			goto err;
		}
		if (len > 255) {
			tls_config_set_errorx(config,
			    "alpn protocol too long");
			goto err;
		}
		buf[i++] = len & 0xff;
		memcpy(&buf[i], p, len);
		i += len;
	}

	free(s);

	*alpn_data = buf;
	*alpn_len = buf_len;

	return (0);

 err:
	free(buf);
	free(s);

	return (-1);
}

int
tls_config_set_alpn(struct tls_config *config, const char *alpn)
{
	return tls_config_parse_alpn(config, alpn, &config->alpn,
	    &config->alpn_len);
}

static int
tls_config_add_keypair_file_internal(struct tls_config *config,
    const char *cert_file, const char *key_file, const char *ocsp_file)
{
	struct tls_keypair *keypair;

	if ((keypair = tls_keypair_new()) == NULL)
		return (-1);
	if (tls_keypair_set_cert_file(keypair, &config->error, cert_file) != 0)
		goto err;
	if (tls_keypair_set_key_file(keypair, &config->error, key_file) != 0)
		goto err;
	if (ocsp_file != NULL &&
	    tls_keypair_set_ocsp_staple_file(keypair, &config->error,
		ocsp_file) != 0)
		goto err;

	tls_config_keypair_add(config, keypair);

	return (0);

 err:
	tls_keypair_free(keypair);
	return (-1);
}

static int
tls_config_add_keypair_mem_internal(struct tls_config *config, const uint8_t *cert,
    size_t cert_len, const uint8_t *key, size_t key_len,
    const uint8_t *staple, size_t staple_len)
{
	struct tls_keypair *keypair;

	if ((keypair = tls_keypair_new()) == NULL)
		return (-1);
	if (tls_keypair_set_cert_mem(keypair, cert, cert_len) != 0)
		goto err;
	if (tls_keypair_set_key_mem(keypair, key, key_len) != 0)
		goto err;
	if (staple != NULL &&
	    tls_keypair_set_ocsp_staple_mem(keypair, staple, staple_len) != 0)
		goto err;

	tls_config_keypair_add(config, keypair);

	return (0);

 err:
	tls_keypair_free(keypair);
	return (-1);
}

int
tls_config_add_keypair_mem(struct tls_config *config, const uint8_t *cert,
    size_t cert_len, const uint8_t *key, size_t key_len)
{
	return tls_config_add_keypair_mem_internal(config, cert, cert_len, key,
	    key_len, NULL, 0);
}

int
tls_config_add_keypair_file(struct tls_config *config,
    const char *cert_file, const char *key_file)
{
	return tls_config_add_keypair_file_internal(config, cert_file,
	    key_file, NULL);
}

int
tls_config_add_keypair_ocsp_mem(struct tls_config *config, const uint8_t *cert,
    size_t cert_len, const uint8_t *key, size_t key_len, const uint8_t *staple,
    size_t staple_len)
{
	return tls_config_add_keypair_mem_internal(config, cert, cert_len, key,
	    key_len, staple, staple_len);
}

int
tls_config_add_keypair_ocsp_file(struct tls_config *config,
    const char *cert_file, const char *key_file, const char *ocsp_file)
{
	return tls_config_add_keypair_file_internal(config, cert_file,
	    key_file, ocsp_file);
}

int
tls_config_set_ca_file(struct tls_config *config, const char *ca_file)
{
	return tls_config_load_file(&config->error, "CA", ca_file,
	    &config->ca_mem, &config->ca_len);
}

int
tls_config_set_ca_path(struct tls_config *config, const char *ca_path)
{
	return set_string(&config->ca_path, ca_path);
}

int
tls_config_set_ca_mem(struct tls_config *config, const uint8_t *ca, size_t len)
{
	return set_mem(&config->ca_mem, &config->ca_len, ca, len);
}

int
tls_config_set_cert_file(struct tls_config *config, const char *cert_file)
{
	return tls_keypair_set_cert_file(config->keypair, &config->error,
	    cert_file);
}

int
tls_config_set_cert_mem(struct tls_config *config, const uint8_t *cert,
    size_t len)
{
	return tls_keypair_set_cert_mem(config->keypair, cert, len);
}

int
tls_config_set_ciphers(struct tls_config *config, const char *ciphers)
{
	SSL_CTX *ssl_ctx = NULL;

	if (ciphers == NULL ||
	    strcasecmp(ciphers, "default") == 0 ||
	    strcasecmp(ciphers, "secure") == 0)
		ciphers = TLS_CIPHERS_DEFAULT;
	else if (strcasecmp(ciphers, "compat") == 0)
		ciphers = TLS_CIPHERS_COMPAT;
	else if (strcasecmp(ciphers, "legacy") == 0)
		ciphers = TLS_CIPHERS_LEGACY;
	else if (strcasecmp(ciphers, "all") == 0 ||
	    strcasecmp(ciphers, "insecure") == 0)
		ciphers = TLS_CIPHERS_ALL;

	if ((ssl_ctx = SSL_CTX_new(SSLv23_method())) == NULL) {
		tls_config_set_errorx(config, "out of memory");
		goto fail;
	}
	if (SSL_CTX_set_cipher_list(ssl_ctx, ciphers) != 1) {
		tls_config_set_errorx(config, "no ciphers for '%s'", ciphers);
		goto fail;
	}

	SSL_CTX_free(ssl_ctx);
	return set_string(&config->ciphers, ciphers);

 fail:
	SSL_CTX_free(ssl_ctx);
	return -1;
}

int
tls_config_set_dheparams(struct tls_config *config, const char *params)
{
	int keylen;

	if (params == NULL || strcasecmp(params, "none") == 0)
		keylen = 0;
	else if (strcasecmp(params, "auto") == 0)
		keylen = -1;
	else if (strcasecmp(params, "legacy") == 0)
		keylen = 1024;
	else {
		tls_config_set_errorx(config, "invalid dhe param '%s'", params);
		return (-1);
	}

	config->dheparams = keylen;

	return (0);
}

int
tls_config_set_ecdhecurve(struct tls_config *config, const char *name)
{
	int nid;

	if (name == NULL || strcasecmp(name, "none") == 0)
		nid = NID_undef;
	else if (strcasecmp(name, "auto") == 0)
		nid = -1;
	else if ((nid = OBJ_txt2nid(name)) == NID_undef) {
		tls_config_set_errorx(config, "invalid ecdhe curve '%s'", name);
		return (-1);
	}

	config->ecdhecurve = nid;

	return (0);
}

int
tls_config_set_key_file(struct tls_config *config, const char *key_file)
{
	return tls_keypair_set_key_file(config->keypair, &config->error,
	    key_file);
}

int
tls_config_set_key_mem(struct tls_config *config, const uint8_t *key,
    size_t len)
{
	return tls_keypair_set_key_mem(config->keypair, key, len);
}

static int
tls_config_set_keypair_file_internal(struct tls_config *config,
    const char *cert_file, const char *key_file, const char *ocsp_file)
{
	if (tls_config_set_cert_file(config, cert_file) != 0)
		return (-1);
	if (tls_config_set_key_file(config, key_file) != 0)
		return (-1);
	if (tls_config_set_key_file(config, key_file) != 0)
		return (-1);
	if (ocsp_file != NULL &&
	    tls_config_set_ocsp_staple_file(config, ocsp_file) != 0)
		return (-1);

	return (0);
}

static int
tls_config_set_keypair_mem_internal(struct tls_config *config, const uint8_t *cert,
    size_t cert_len, const uint8_t *key, size_t key_len,
    const uint8_t *staple, size_t staple_len)
{
	if (tls_config_set_cert_mem(config, cert, cert_len) != 0)
		return (-1);
	if (tls_config_set_key_mem(config, key, key_len) != 0)
		return (-1);
	if ((staple != NULL) &&
	    (tls_config_set_ocsp_staple_mem(config, staple, staple_len) != 0))
		return (-1);

	return (0);
}

int
tls_config_set_keypair_file(struct tls_config *config,
    const char *cert_file, const char *key_file)
{
	return tls_config_set_keypair_file_internal(config, cert_file, key_file,
	    NULL);
}

int
tls_config_set_keypair_mem(struct tls_config *config, const uint8_t *cert,
    size_t cert_len, const uint8_t *key, size_t key_len)
{
	return tls_config_set_keypair_mem_internal(config, cert, cert_len,
	    key, key_len, NULL, 0);
}

int
tls_config_set_keypair_ocsp_file(struct tls_config *config,
    const char *cert_file, const char *key_file, const char *ocsp_file)
{
	return tls_config_set_keypair_file_internal(config, cert_file, key_file,
	    ocsp_file);
}

int
tls_config_set_keypair_ocsp_mem(struct tls_config *config, const uint8_t *cert,
    size_t cert_len, const uint8_t *key, size_t key_len,
    const uint8_t *staple, size_t staple_len)
{
	return tls_config_set_keypair_mem_internal(config, cert, cert_len,
	    key, key_len, staple, staple_len);
}


int
tls_config_set_protocols(struct tls_config *config, uint32_t protocols)
{
	config->protocols = protocols;

	return (0);
}

int
tls_config_set_verify_depth(struct tls_config *config, int verify_depth)
{
	config->verify_depth = verify_depth;

	return (0);
}

void
tls_config_prefer_ciphers_client(struct tls_config *config)
{
	config->ciphers_server = 0;
}

void
tls_config_prefer_ciphers_server(struct tls_config *config)
{
	config->ciphers_server = 1;
}

void
tls_config_insecure_noverifycert(struct tls_config *config)
{
	config->verify_cert = 0;
}

void
tls_config_insecure_noverifyname(struct tls_config *config)
{
	config->verify_name = 0;
}

void
tls_config_insecure_noverifytime(struct tls_config *config)
{
	config->verify_time = 0;
}

void
tls_config_verify(struct tls_config *config)
{
	config->verify_cert = 1;
	config->verify_name = 1;
	config->verify_time = 1;
}

void
tls_config_ocsp_require_stapling(struct tls_config *config)
{
	config->ocsp_require_stapling = 1;
}

void
tls_config_verify_client(struct tls_config *config)
{
	config->verify_client = 1;
}

void
tls_config_verify_client_optional(struct tls_config *config)
{
	config->verify_client = 2;
}

int
tls_config_set_ocsp_staple_file(struct tls_config *config, const char *staple_file)
{
	return tls_keypair_set_ocsp_staple_file(config->keypair, &config->error,
	    staple_file);
}

int
tls_config_set_ocsp_staple_mem(struct tls_config *config, const uint8_t *staple,
    size_t len)
{
	return tls_keypair_set_ocsp_staple_mem(config->keypair, staple, len);
}

int
tls_config_set_session_id(struct tls_config *config,
    const unsigned char *session_id, size_t len)
{
	if (len > TLS_MAX_SESSION_ID_LENGTH) {
		tls_config_set_errorx(config, "session ID too large");
		return (-1);
	}
	memset(config->session_id, 0, sizeof(config->session_id));
	memcpy(config->session_id, session_id, len);
	return (0);
}

int
tls_config_set_session_lifetime(struct tls_config *config, int lifetime)
{
	if (lifetime > TLS_MAX_SESSION_TIMEOUT) {
		tls_config_set_errorx(config, "session lifetime too large");
		return (-1);
	}
	if (lifetime != 0 && lifetime < TLS_MIN_SESSION_TIMEOUT) {
		tls_config_set_errorx(config, "session lifetime too small");
		return (-1);
	}

	config->session_lifetime = lifetime;
	return (0);
}

int
tls_config_add_ticket_key(struct tls_config *config, uint32_t keyrev,
    unsigned char *key, size_t keylen)
{
	struct tls_ticket_key newkey;
	int i;

	if (TLS_TICKET_KEY_SIZE != keylen ||
	    sizeof(newkey.aes_key) + sizeof(newkey.hmac_key) > keylen) {
		tls_config_set_errorx(config,
		    "wrong amount of ticket key data");
		return (-1);
	}

	keyrev = htonl(keyrev);
	memset(&newkey, 0, sizeof(newkey));
	memcpy(newkey.key_name, &keyrev, sizeof(keyrev));
	memcpy(newkey.aes_key, key, sizeof(newkey.aes_key));
	memcpy(newkey.hmac_key, key + sizeof(newkey.aes_key),
	    sizeof(newkey.hmac_key));
	newkey.time = time(NULL);

	for (i = 0; i < TLS_NUM_TICKETS; i++) {
		struct tls_ticket_key *tk = &config->ticket_keys[i];
		if (memcmp(newkey.key_name, tk->key_name,
		    sizeof(tk->key_name)) != 0)
			continue;

		/* allow re-entry of most recent key */
		if (i == 0 && memcmp(newkey.aes_key, tk->aes_key,
		    sizeof(tk->aes_key)) == 0 && memcmp(newkey.hmac_key,
		    tk->hmac_key, sizeof(tk->hmac_key)) == 0)
			return (0);
		tls_config_set_errorx(config, "ticket key already present");
		return (-1);
	}

	memmove(&config->ticket_keys[1], &config->ticket_keys[0],
	    sizeof(config->ticket_keys) - sizeof(config->ticket_keys[0]));
	config->ticket_keys[0] = newkey;

	config->ticket_autorekey = 0;

	return (0);
}

int
tls_config_ticket_autorekey(struct tls_config *config)
{
	unsigned char key[TLS_TICKET_KEY_SIZE];
	int rv;

	arc4random_buf(key, sizeof(key));
	rv = tls_config_add_ticket_key(config, config->ticket_keyrev++, key,
	    sizeof(key));
	config->ticket_autorekey = 1;
	return (rv);
}