document.getElementById("btn-login").addEventListener("click", async () => {
  let options;
  let credential;

  try {
    // passkey 認証のためのオプション取得
    options = await loginOptionAPI();
    console.info("認証オプションの取得完了", { options });

    // 認証
    const publicKeyOption = PublicKeyCredential.parseRequestOptionsFromJSON(
      options.publicKey
    );
    credential = await navigator.credentials.get({
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

  try {
    const token = await loginAPI(credential);

    alert("認証完了！", { token });

    localStorage.token = token;
    window.location.href = "/dashboard.html";
  } catch (err) {
    if (err.name === "NotAllowedError") {
      alert("中断しました");
      return;
    }

    // Password Manager に登録された鍵の削除
    if (PublicKeyCredential.signalUnknownCredential) {
      await PublicKeyCredential.signalUnknownCredential({
        rpId: options.publicKey.rpId,
        credentialId: credential.id,
      });
    }

    console.error(err);
    alert(err.message);
  }
});

/**
 * @returns {Credential}
 */
async function loginOptionAPI() {
  const res = await fetch("http://localhost:8080/login_options", {
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
 * @param {Credential}
 * @returns {string} token
 */
async function loginAPI(credential) {
  const res = await fetch("http://localhost:8080/login", {
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
