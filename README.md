# tiktok
TikTok Web Signatures Reverse Engineered

This repo contains the following algorithms for TikTok's Web SDK:
 - **X-Bogus** - URL signer
 - **Edata** - Captcha data
 - **X-Gnarly** - URL signer (with Bogus)
 - **X-Mssdk** - Fingerprint for certain actions (register/login)
 - **Strdata** - Fingerprint for refreshing **MsToken**

# Usage

- Run `go build .` and it will output a binary called `tiktok`
- Run the `tiktok` binary (the main API server)
- Send requests to the hosted endpoints to encrypt your data

# Notice

- This is for educational and research purposes only
- Do not use this to abuse TikTok's anti-abuse systems
- Some algorithms may have changed- this is SDK version `5.1.1` (current is `5.1.3-ZTCA`)

# Contact

Contact via `Telegram` or `email` is available

# Takedowns
⚠️I am **fully willing to comply** with any takedown request by the order of *TikTok*, *ByteDance*, or *associated parties* ⚠️
