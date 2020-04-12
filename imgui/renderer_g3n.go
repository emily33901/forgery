package imgui

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/window"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/inkyblackness/imgui-go"
)

// Opengl3 Implements rendering with g3n as a backend
// G3nOpenGL3 implements a renderer based on github.com/go-gl/gl (v3.2-core).
type G3nOpenGL3 struct {
	imguiIO imgui.IO

	gl      *gls.GLS
	program *gls.Program

	glslVersion            string
	fontTexture            uint32
	vertHandle             uint32
	fragHandle             uint32
	attribLocationTex      int32
	attribLocationProjMtx  int32
	attribLocationPosition int32
	attribLocationUV       int32
	attribLocationColor    int32
	vboHandle              uint32
	elementsHandle         uint32
}

// NewG3nOpenGL3 attempts to initialize a renderer.
// An OpenGL context has to be established before calling this function.
func NewG3nOpenGL3(io imgui.IO) (*G3nOpenGL3, error) {
	renderer := &G3nOpenGL3{
		imguiIO:     io,
		gl:          window.Get().Gls(),
		glslVersion: "#version 330",
	}

	renderer.createDeviceObjects()
	return renderer, nil
}

// Dispose cleans up the resources.
func (renderer *G3nOpenGL3) Dispose() {
	renderer.invalidateDeviceObjects()
}

// PreRender clears the framebuffer.
func (renderer *G3nOpenGL3) PreRender(clearColor [3]float32) {
	renderer.gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], 1.0)
	renderer.gl.Clear(gl.COLOR_BUFFER_BIT)
}

// Render translates the ImGui draw data to OpenGL3 commands.
func (renderer *G3nOpenGL3) Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData) {
	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	displayWidth, displayHeight := displaySize[0], displaySize[1]
	fbWidth, fbHeight := framebufferSize[0], framebufferSize[1]
	if (fbWidth <= 0) || (fbHeight <= 0) {
		return
	}
	drawData.ScaleClipRects(imgui.Vec2{
		X: fbWidth / displayWidth,
		Y: fbHeight / displayHeight,
	})

	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(gl.TEXTURE0)
	var lastProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &lastProgram)
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	var lastSampler int32
	gl.GetIntegerv(gl.SAMPLER_BINDING, &lastSampler)
	var lastArrayBuffer int32
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	var lastElementArrayBuffer int32
	gl.GetIntegerv(gl.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	var lastVertexArray int32
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(gl.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	var lastBlendSrcRgb int32
	gl.GetIntegerv(gl.BLEND_SRC_RGB, &lastBlendSrcRgb)
	var lastBlendDstRgb int32
	gl.GetIntegerv(gl.BLEND_DST_RGB, &lastBlendDstRgb)
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(gl.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	var lastBlendDstAlpha int32
	gl.GetIntegerv(gl.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	var lastBlendEquationRgb int32
	gl.GetIntegerv(gl.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(gl.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	lastEnableBlend := gl.IsEnabled(gl.BLEND)
	lastEnableCullFace := gl.IsEnabled(gl.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(gl.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(gl.SCISSOR_TEST)

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	renderer.gl.Enable(gl.BLEND)
	renderer.gl.BlendEquation(gl.FUNC_ADD)
	renderer.gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	renderer.gl.Disable(gl.CULL_FACE)
	renderer.gl.Disable(gl.DEPTH_TEST)
	renderer.gl.Enable(gl.SCISSOR_TEST)
	renderer.gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// Setup viewport, orthographic projection matrix
	// Our visible imgui space lies from draw_data->DisplayPos (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right).
	// DisplayMin is typically (0,0) for single viewport apps.
	renderer.gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	orthoProjection := [4][4]float32{
		{2.0 / displayWidth, 0.0, 0.0, 0.0},
		{0.0, 2.0 / -displayHeight, 0.0, 0.0},
		{0.0, 0.0, -1.0, 0.0},
		{-1.0, 1.0, 0.0, 1.0},
	}
	renderer.gl.UseProgram(renderer.program)
	renderer.gl.Uniform1i(renderer.attribLocationTex, 0)
	renderer.gl.UniformMatrix4fv(renderer.attribLocationProjMtx, 1, false, &orthoProjection[0][0])
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and
	// we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	vaoHandle := renderer.gl.GenVertexArray()
	renderer.gl.BindVertexArray(vaoHandle)
	renderer.gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vboHandle)
	renderer.gl.EnableVertexAttribArray(uint32(renderer.attribLocationPosition))
	renderer.gl.EnableVertexAttribArray(uint32(renderer.attribLocationUV))
	renderer.gl.EnableVertexAttribArray(uint32(renderer.attribLocationColor))
	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	renderer.gl.VertexAttribPointer(uint32(renderer.attribLocationPosition), 2, gl.FLOAT, false, int32(vertexSize), uint32(vertexOffsetPos))
	renderer.gl.VertexAttribPointer(uint32(renderer.attribLocationUV), 2, gl.FLOAT, false, int32(vertexSize), uint32(vertexOffsetUv))
	renderer.gl.VertexAttribPointer(uint32(renderer.attribLocationColor), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), uint32(vertexOffsetCol))
	indexSize := imgui.IndexBufferLayout()
	drawType := gl.UNSIGNED_SHORT
	if indexSize == 4 {
		drawType = gl.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {
		var indexBufferOffset uint32

		vBuffer, vertexBufferSize := list.VertexBuffer()
		vertexBuffer := (*uint8)(vBuffer)
		renderer.gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vboHandle)
		renderer.gl.BufferData(gl.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, gl.STREAM_DRAW)

		iBuffer, indexBufferSize := list.IndexBuffer()
		indexBuffer := (*uint8)(iBuffer)
		renderer.gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, renderer.elementsHandle)
		renderer.gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, gl.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {
				renderer.gl.BindTexture(gl.TEXTURE_2D, uint32(cmd.TextureID()))
				clipRect := cmd.ClipRect()
				renderer.gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), uint32(clipRect.Z-clipRect.X), uint32(clipRect.W-clipRect.Y))
				renderer.gl.DrawElements(gl.TRIANGLES, int32(cmd.ElementCount()), uint32(drawType), indexBufferOffset)
			}
			indexBufferOffset += uint32(cmd.ElementCount() * indexSize)
		}
	}
	renderer.gl.DeleteVertexArrays(vaoHandle)

	// Restore modified GL state
	gl.UseProgram(uint32(lastProgram))
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindSampler(0, uint32(lastSampler))
	gl.ActiveTexture(uint32(lastActiveTexture))
	gl.BindVertexArray(uint32(lastVertexArray))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	if lastEnableBlend {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	if lastEnableCullFace {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
	if lastEnableDepthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	if lastEnableScissorTest {
		gl.Enable(gl.SCISSOR_TEST)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}
	gl.PolygonMode(gl.FRONT_AND_BACK, uint32(lastPolygonMode[0]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}

func (renderer *G3nOpenGL3) createDeviceObjects() {
	gl.Init()

	vertexShader := renderer.glslVersion + `
uniform mat4 ProjMtx;
in vec2 Position;
in vec2 UV;
in vec4 Color;
out vec2 Frag_UV;
out vec4 Frag_Color;
void main()
{
	Frag_UV = UV;
	Frag_Color = Color;
	gl_Position = ProjMtx * vec4(Position.xy,0,1);
}
`
	fragmentShader := renderer.glslVersion + `
uniform sampler2D Texture;
in vec2 Frag_UV;
in vec4 Frag_Color;
out vec4 Out_Color;
void main()
{
	Out_Color = vec4(Frag_Color.rgb, Frag_Color.a * texture( Texture, Frag_UV.st).r);
}
`

	prog := renderer.gl.NewProgram()
	prog.AddShader(gls.VERTEX_SHADER, vertexShader)
	prog.AddShader(gls.FRAGMENT_SHADER, fragmentShader)

	err := prog.Build()

	if err != nil {
		panic(err)
	}

	renderer.program = prog
	renderer.attribLocationTex = prog.GetUniformLocation("Texture")
	renderer.attribLocationProjMtx = prog.GetUniformLocation("ProjMtx")
	renderer.attribLocationPosition = prog.GetAttribLocation("Position")
	renderer.attribLocationUV = prog.GetAttribLocation("UV")
	renderer.attribLocationColor = prog.GetAttribLocation("Color")

	renderer.vboHandle = renderer.gl.GenBuffer()
	renderer.elementsHandle = renderer.gl.GenBuffer()

	renderer.createFontsTexture()
}

func (renderer *G3nOpenGL3) createFontsTexture() {
	// Build texture atlas
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataAlpha8()

	renderer.fontTexture = renderer.gl.GenTexture()
	renderer.gl.BindTexture(gls.TEXTURE_2D, renderer.fontTexture)
	renderer.gl.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MIN_FILTER, gls.LINEAR)
	renderer.gl.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MAG_FILTER, gls.LINEAR)
	renderer.gl.TexImage2D(gls.TEXTURE_2D, 0, gls.RED, int32(image.Width), int32(image.Height), gls.RED, gl.UNSIGNED_BYTE, (*uint8)(image.Pixels))

	io.Fonts().SetTextureID(imgui.TextureID(renderer.fontTexture))
}

func (renderer *G3nOpenGL3) invalidateDeviceObjects() {
	if renderer.vboHandle != 0 {
		renderer.gl.DeleteBuffers(renderer.vboHandle)
	}
	renderer.vboHandle = 0
	if renderer.elementsHandle != 0 {
		renderer.gl.DeleteBuffers(renderer.elementsHandle)
	}
	renderer.elementsHandle = 0

	renderer.gl.DeleteProgram(renderer.program.Handle())
	renderer.program = nil

	if renderer.fontTexture != 0 {
		renderer.gl.DeleteTextures(renderer.fontTexture)
		gl.DeleteTextures(1, &renderer.fontTexture)
		imgui.CurrentIO().Fonts().SetTextureID(0)
		renderer.fontTexture = 0
	}
}
