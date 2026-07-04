# Vibing and AI

A curated collection of AI-assisted projects for experimentation, testing, and educational purposes.

## 📋 Purpose

This repository serves as a sandbox and learning environment for exploring AI-assisted development workflows, tools, and methodologies. Each project within this repository demonstrates different aspects of AI collaboration in software development, from proof-of-concepts to educational examples.

## 🏗️ Repository Structure

Each project is **standalone** and contained within its own subdirectory, ensuring:
- **Isolation**: Projects don't interfere with each other
- **Clarity**: Each project has its own dependencies and documentation
- **Portability**: Projects can be easily extracted or shared independently

```
vibing-and-ai/
├── project-1/
│   ├── README.md
│   ├── src/
│   └── ...
├── project-2/
│   ├── README.md
│   ├── src/
│   └── ...
└── ...
```

## 🎯 Use Cases

This repository is ideal for:

- **Experimentation**: Testing new AI tools, frameworks, and development approaches
- **Learning**: Educational projects demonstrating AI-assisted coding techniques
- **Prototyping**: Quick proof-of-concepts and MVPs built with AI assistance
- **Benchmarking**: Comparing different AI coding assistants and methodologies
- **Documentation**: Showcasing best practices for human-AI collaboration

## 🚀 Getting Started

1. **Browse Projects**: Navigate to individual project folders to explore different experiments
2. **Read Documentation**: Each project contains its own README with specific setup instructions
3. **Run Independently**: Follow the instructions in each project's directory to run it standalone
4. **Learn and Adapt**: Use these projects as templates or learning resources for your own work

## 📁 Project Guidelines

Each project should include:

- **README.md**: Clear documentation of purpose, setup, and usage
- **Dependencies**: Explicit listing of requirements (package.json, requirements.txt, etc.)
- **Self-contained**: All necessary files and resources within the project folder
- **AI Context**: Optional notes on which AI tools were used and how

## 🤖 Issue Management for AI Agents

This repository follows an **AI-first issue tracking approach**. All issues should be written with AI coding agents in mind, enabling them to understand and resolve tasks independently.

### Issue Requirements

Every issue must include:

1. **Clear Title**: Concise description of the feature or bug (e.g., "Add user authentication to project-auth" or "Fix memory leak in project-data-viz")

2. **Detailed Description**: Comprehensive context including:
   - **Current State**: What exists now
   - **Desired State**: What should exist after completion
   - **Acceptance Criteria**: Specific, testable conditions for completion
   - **Technical Context**: Relevant technologies, frameworks, or constraints
   - **File Locations**: Specific files or directories that need modification

3. **Project Label**: Every issue must be tagged with the corresponding project subdirectory
   - Use labels like `project:project-name` (e.g., `project:auth-system`, `project:data-viz`)
   - This ensures AI agents can identify which project the issue belongs to

4. **AI-Friendly Format**: Structure issues for machine readability:
   ```markdown
   ## Context
   [Background information]
   
   ## Problem/Feature
   [Detailed description]
   
   ## Expected Outcome
   [What success looks like]
   
   ## Technical Requirements
   - Requirement 1
   - Requirement 2
   
   ## Files to Modify
   - `path/to/file1.js`
   - `path/to/file2.py`
   
   ## Testing Criteria
   - [ ] Test case 1
   - [ ] Test case 2
   ```

### Issue Template Example

```markdown
Title: Implement dark mode toggle in project-ui-components

Labels: project:ui-components, enhancement, good-first-issue

## Context
The UI components project currently only supports light mode. Users have requested a dark mode option.

## Problem/Feature
Add a dark mode toggle that switches between light and dark color schemes across all components.

## Expected Outcome
- A toggle button in the header
- Persistent theme preference (localStorage)
- All components render correctly in both modes
- Smooth transition animations between modes

## Technical Requirements
- Use CSS custom properties for theming
- Implement React Context for theme state management
- Add toggle component to Header.jsx
- Update all component stylesheets to support both themes

## Files to Modify
- `project-ui-components/src/components/Header.jsx`
- `project-ui-components/src/context/ThemeContext.jsx` (new file)
- `project-ui-components/src/styles/themes.css` (new file)
- `project-ui-components/src/components/**/*.css` (update all)

## Testing Criteria
- [ ] Toggle switches between light and dark modes
- [ ] Theme preference persists across page refreshes
- [ ] All components display correctly in both modes
- [ ] Transitions are smooth (300ms)
- [ ] No console errors or warnings
```

### Benefits of AI-First Issues

- **Autonomous Resolution**: AI agents can pick up and complete issues without human intervention
- **Consistency**: Standardized format ensures predictable AI performance
- **Traceability**: Clear project labels enable easy filtering and organization
- **Quality**: Detailed requirements reduce ambiguity and implementation errors

## 🤝 Contributing

Feel free to add new experimental projects following the established structure. Each contribution should:

1. Be placed in its own subdirectory
2. Include comprehensive documentation
3. Be fully functional and self-contained
4. Optionally document the AI assistance methodology used

## 📝 License

Individual projects may have their own licenses. Please refer to each project's directory for specific licensing information.

---

**Note**: This is an experimental repository. Projects may be in various states of completion and are primarily for learning and exploration purposes.
