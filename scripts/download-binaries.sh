#!/bin/bash

set -e

BINARIES_DIR="internal/bootstrap/binaries"
MAX_PARALLEL=6
PIDS=()

mkdir -p ${BINARIES_DIR}/{darwin_amd64,darwin_arm64,linux_amd64,linux_arm64,windows_amd64,windows_arm64}

echo "Downloading yt-dlp and ffmpeg binaries for all platforms..."

download_with_retry() {
    local url=$1
    local dest=$2
    local max_attempts=3
    local attempt=1

    if [ -f "$dest" ] && [ -s "$dest" ]; then
        echo "  ✓ Skipping $dest (already exists)"
        return 0
    fi

    while [ $attempt -le $max_attempts ]; do
        if curl -L -f --connect-timeout 10 --max-time 300 -o "$dest" "$url" 2>/dev/null; then
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 1
    done
    return 1
}

wait_for_slot() {
    while [ ${#PIDS[@]} -ge $MAX_PARALLEL ]; do
        for pid in "${!PIDS[@]}"; do
            if ! kill -0 "${PIDS[$pid]}" 2>/dev/null; then
                unset PIDS[$pid]
            fi
        done
        sleep 0.1
    done
}

download_async() {
    local url=$1
    local dest=$2
    local name=$3
    
    wait_for_slot
    
    (
        if download_with_retry "$url" "$dest"; then
            echo "  ✓ Downloaded $name"
        else
            echo "  ✗ Failed to download $name" >&2
        fi
    ) &
    
    PIDS+=($!)
}

echo "Downloading yt-dlp..."
download_async "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos" "${BINARIES_DIR}/darwin_amd64/yt-dlp" "yt-dlp (darwin_amd64)"
download_async "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos_legacy" "${BINARIES_DIR}/darwin_arm64/yt-dlp" "yt-dlp (darwin_arm64)" || \
download_async "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos" "${BINARIES_DIR}/darwin_arm64/yt-dlp" "yt-dlp (darwin_arm64)"
download_async "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux" "${BINARIES_DIR}/linux_amd64/yt-dlp" "yt-dlp (linux_amd64)"
download_async "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux_aarch64" "${BINARIES_DIR}/linux_arm64/yt-dlp" "yt-dlp (linux_arm64)"
download_async "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" "${BINARIES_DIR}/windows_amd64/yt-dlp.exe" "yt-dlp (windows_amd64)"
download_async "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" "${BINARIES_DIR}/windows_arm64/yt-dlp.exe" "yt-dlp (windows_arm64)"

for pid in "${PIDS[@]}"; do
    wait $pid 2>/dev/null || true
done
PIDS=()

chmod +x ${BINARIES_DIR}/*/yt-dlp* 2>/dev/null || true

echo "Downloading ffmpeg..."

download_and_extract_ffmpeg() {
    local url=$1
    local platform=$2
    local archive_ext=$3
    local binary_name=$4
    
    wait_for_slot
    
    (
        local archive_file="${BINARIES_DIR}/${platform}/ffmpeg.${archive_ext}"
        
        if [ -f "${BINARIES_DIR}/${platform}/${binary_name}" ] && [ -s "${BINARIES_DIR}/${platform}/${binary_name}" ]; then
            echo "  ✓ Skipping ffmpeg for ${platform} (already exists)"
            return 0
        fi
        
        if download_with_retry "$url" "$archive_file"; then
            case "$archive_ext" in
                zip)
                    unzip -q -j -o "$archive_file" -d "${BINARIES_DIR}/${platform}/" 2>/dev/null || \
                    unzip -q -o "$archive_file" -d "${BINARIES_DIR}/${platform}/" 2>/dev/null || true
                    find "${BINARIES_DIR}/${platform}/" -name "$binary_name" -type f -exec mv {} "${BINARIES_DIR}/${platform}/${binary_name}" \; 2>/dev/null || true
                    rm -f "$archive_file" "${BINARIES_DIR}/${platform}/"*.dylib 2>/dev/null || true
                    rm -rf "${BINARIES_DIR}/${platform}/ffmpeg-"* 2>/dev/null || true
                    ;;
                tar.xz)
                    tar -xf "$archive_file" -C "${BINARIES_DIR}/${platform}/" --strip-components=1 --wildcards "*/ffmpeg" 2>/dev/null || true
                    rm -f "$archive_file" "${BINARIES_DIR}/${platform}/"*.md "${BINARIES_DIR}/${platform}/"*.txt 2>/dev/null || true
                    ;;
            esac
            echo "  ✓ Downloaded and extracted ffmpeg for ${platform}"
        else
            echo "  ✗ Failed to download ffmpeg for ${platform}" >&2
        fi
    ) &
    
    PIDS+=($!)
}

download_and_extract_ffmpeg "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip" "darwin_amd64" "zip" "ffmpeg"
download_and_extract_ffmpeg "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip" "darwin_arm64" "zip" "ffmpeg"
download_and_extract_ffmpeg "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz" "linux_amd64" "tar.xz" "ffmpeg"
download_and_extract_ffmpeg "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz" "linux_arm64" "tar.xz" "ffmpeg"
download_and_extract_ffmpeg "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" "windows_amd64" "zip" "ffmpeg.exe"
download_and_extract_ffmpeg "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" "windows_arm64" "zip" "ffmpeg.exe"

for pid in "${PIDS[@]}"; do
    wait $pid 2>/dev/null || true
done

chmod +x ${BINARIES_DIR}/*/ffmpeg* 2>/dev/null || true

echo ""
echo "Download complete!"
echo "Verifying binaries..."
for platform in darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 windows_amd64 windows_arm64; do
    ytdlp_ok=false
    ffmpeg_ok=false
    
    if [ -f "${BINARIES_DIR}/${platform}/yt-dlp" ] || [ -f "${BINARIES_DIR}/${platform}/yt-dlp.exe" ]; then
        ytdlp_ok=true
        echo "✓ yt-dlp found for ${platform}"
    else
        echo "✗ yt-dlp missing for ${platform}"
    fi
    
    if [ -f "${BINARIES_DIR}/${platform}/ffmpeg" ] || [ -f "${BINARIES_DIR}/${platform}/ffmpeg.exe" ]; then
        ffmpeg_ok=true
        echo "✓ ffmpeg found for ${platform}"
    else
        echo "✗ ffmpeg missing for ${platform}"
    fi
done
