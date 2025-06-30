const usernameEl = document.getElementById("username");

document.getElementById("btn-register").addEventListener("click", async () => {
  let options;
  let credential;

  try {
    // ユーザーの登録
    await signUpAPI(usernameEl.value);
    console.info("ユーザー登録完了");

    // passkey 登録のためのオプション取得
    options = await registerOptionsAPI();
    console.info("登録オプションの取得完了", { options });

    // 認証
    const publicKeyOption = PublicKeyCredential.parseCreationOptionsFromJSON(
      options.publicKey
    );
    credential = await navigator.credentials.create({
      publicKey: publicKeyOption,
    });
    console.info("認証成功", { credential });
  } catch (err) {
    if (err.name === "NotAllowedError") {
      alert("中断しました");
      return;
    }
    console.error(err);
    alert(err.message);

    return;
  }

  // 公開鍵の登録
  try {
    const token = await registerAPI(credential.toJSON());

    localStorage.token = token;
    console.info("登録完了", { token });

    alert("登録完了！");
    window.location.href = "/dashboard.html";
  } catch (err) {
    if (err.name === "NotAllowedError") {
      alert("中断しました");
      return;
    }

    // Password Manager に登録された鍵の削除
    if (PublicKeyCredential.signalUnknownCredential) {
      await PublicKeyCredential.signalUnknownCredential({
        rpId: options.publicKey.rp.id,
        credentialId: credential.id,
      });
    }

    console.error(err);
    alert(err.message);
  }
});

/**
 * @param {string} username
 */
async function signUpAPI(username) {
  const res = await fetch("http://localhost:8080/signup", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username }),
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data?.message ?? res.statusText);
  }
}

/**
 * @returns {PublicKeyCredentialCreationOptions}
 */
async function registerOptionsAPI() {
  const res = await fetch("http://localhost:8080/register_options", {
    method: "GET",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data?.message ?? res.statusText);
  }

  return data.options;
}

/**
 * @param {Credential} credential
 * @returns {string} token
 */
async function registerAPI(credential) {
  const res = await fetch("http://localhost:8080/register", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credential),
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data?.message ?? res.statusText);
  }

  return data.token;
}
