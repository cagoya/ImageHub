// 图片处理 Web Worker

/**
 * RGB转HSL (返回 0-1 范围)
 */
function rgbToHsl(r: number, g: number, b: number): [number, number, number] {
  r /= 255;
  g /= 255;
  b /= 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h = 0, s = 0, l = (max + min) / 2;

  if (max !== min) {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
    h /= 6;
  }

  return [h, s, l];
}

/**
 * HSL转RGB
 */
function hslToRgb(h: number, s: number, l: number): [number, number, number] {
  const hue2rgb = (p: number, q: number, t: number) => {
    if (t < 0) t += 1;
    if (t > 1) t -= 1;
    if (t < 1/6) return p + (q - p) * 6 * t;
    if (t < 1/2) return q;
    if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
    return p;
  };

  let r, g, b;

  if (s === 0) {
    r = g = b = l;
  } else {
    const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    const p = 2 * l - q;
    r = hue2rgb(p, q, h + 1/3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1/3);
  }

  return [r * 255, g * 255, b * 255];
}

/**
 * 处理图片数据
 */
function processImageData(
  imageData: ImageData,
  brightness: number,
  contrast: number,
  saturation: number,
  hue: number
): ImageData {
  const data = imageData.data;
  const length = data.length;

  // 计算调整系数
  const brightnessValue = brightness * 2.55; // 转换到 -255 到 255
  const contrastFactor = (contrast + 100) / 100;
  const saturationFactor = (saturation + 100) / 100;
  const hueRotation = hue;

  // 处理每个像素
  for (let i = 0; i < length; i += 4) {
    let r = data[i] ?? 0;
    let g = data[i + 1] ?? 0;
    let b = data[i + 2] ?? 0;

    // 应用亮度
    r += brightnessValue;
    g += brightnessValue;
    b += brightnessValue;

    // 应用对比度
    r = ((r / 255 - 0.5) * contrastFactor + 0.5) * 255;
    g = ((g / 255 - 0.5) * contrastFactor + 0.5) * 255;
    b = ((b / 255 - 0.5) * contrastFactor + 0.5) * 255;

    // 转换到 HSL 来应用饱和度和色相
    if (saturationFactor !== 1 || hueRotation !== 0) {
      const hsl = rgbToHsl(r, g, b);
      
      // 应用饱和度
      hsl[1] = Math.max(0, Math.min(1, hsl[1] * saturationFactor));
      
      // 应用色相旋转
      if (hueRotation !== 0) {
        hsl[0] = (hsl[0] + hueRotation / 360) % 1;
        if (hsl[0] < 0) hsl[0] += 1;
      }
      
      // 转换回 RGB
      const rgb = hslToRgb(hsl[0], hsl[1], hsl[2]);
      r = rgb[0];
      g = rgb[1];
      b = rgb[2];
    }

    // 限制在有效范围内
    data[i] = Math.max(0, Math.min(255, r));
    data[i + 1] = Math.max(0, Math.min(255, g));
    data[i + 2] = Math.max(0, Math.min(255, b));
    // data[i + 3] 是 alpha 通道，保持不变
  }

  return imageData;
}

// 监听主线程消息
self.addEventListener('message', (e: MessageEvent) => {
  const { imageData, adjustments, taskId } = e.data;

  try {
    // 重新创建 ImageData 对象
    const imageDataObj = new ImageData(
      new Uint8ClampedArray(imageData.data),
      imageData.width,
      imageData.height
    );

    // 处理图片数据
    const processedImageData = processImageData(
      imageDataObj,
      adjustments.brightness,
      adjustments.contrast,
      adjustments.saturation,
      adjustments.hue
    );

    // 将处理后的数据发送回主线程
    // 使用 Transferable 对象提高性能
    const dataArray = new Uint8ClampedArray(processedImageData.data);
    (self as any).postMessage({
      taskId,
      imageData: {
        data: dataArray,
        width: processedImageData.width,
        height: processedImageData.height
      },
      success: true
    }, [dataArray.buffer]);
  } catch (error) {
    // 发送错误信息
    self.postMessage({
      taskId,
      error: error instanceof Error ? error.message : 'Unknown error',
      success: false
    });
  }
});

