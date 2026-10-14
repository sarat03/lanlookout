import { JSX } from "solid-js";

// One settings row: label + help on the left, control on the right (stacks on phones).
function Field(_props: { id?: string, label: string, help?: JSX.Element, tight?: boolean, children: JSX.Element }) {
  return (
    <div class="wy-field">
      <div>
        {_props.id ? <label for={_props.id}>{_props.label}</label> : <div class="t">{_props.label}</div>}
        <div class="help">{_props.help}</div>
      </div>
      <div class={_props.tight ? "ctl tight" : "ctl"}>{_props.children}</div>
    </div>
  )
}

export default Field
