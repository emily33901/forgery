package render

type FrameBuffer struct {
	adapter      Adapter
	framebuffer  uint32
	colorTexture uint32
	depthTexture uint32
	width        int
	height       int
}

// Resize resizes this framebuffer object
func (win *FrameBuffer) Resize(width int, height int) {
	// if win.width > width && win.height > height {
	// 	// already big enough
	// 	return
	// }

	// round the width and height up to a multiple of 2
	if width%2 == 1 {
		win.width = width + 1
	} else {
		win.width = width
	}

	if height%2 == 1 {
		win.height = height + 1
	} else {
		win.height = height
	}

	win.Bind()

	if win.colorTexture != 0 {
		win.adapter.DeleteTextures(1, &win.colorTexture)
		win.adapter.DeleteRenderBuffer(1, &win.depthTexture)
	}

	win.depthTexture = win.adapter.CreateRenderbufferStorageDepth(int32(win.width), int32(win.height))

	win.adapter.CreateTextureStorage2D(&win.colorTexture, int32(win.width), int32(win.height))
	win.adapter.BindTexture2D(win.colorTexture)
	win.adapter.BindTexture2DToFramebuffer(win.colorTexture)
	win.adapter.BindDepthBufferToFramebuffer(win.depthTexture)
	win.adapter.DrawBuffers()
	win.adapter.ClearColor(0, 0, 0, 0)
	win.adapter.ClearAll()
	win.adapter.BindTexture2D(0)

	win.Unbind()
}

// Bind this framebuffer
func (f *FrameBuffer) Bind() {
	f.adapter.BindFramebuffer(f.framebuffer)
}

// Unbind unbind this framebuffer
func (f *FrameBuffer) Unbind() {
	f.adapter.BindFramebuffer(0)
}

// Destroy deletes and cleans up this framebuffer
func (f *FrameBuffer) Destroy() {
	f.adapter.DeleteFramebuffers(1, &f.framebuffer)
}

func (f *FrameBuffer) TextureId() uint32 {
	return f.colorTexture
}

func (f *FrameBuffer) Size() (int, int) {
	return f.width, f.height
}

// NewFbo returns a new framebuffer
func NewFramebuffer(adapter Adapter, width int, height int) *FrameBuffer {
	f := &FrameBuffer{
		adapter: adapter,
	}
	f.adapter.CreateFramebuffers(1, &f.framebuffer)
	f.Resize(width, height)
	return f
}
