import { setShow, show } from "../functions/exports";

function HistShow(_props: any) {

  const handleSaveShow = (showStr: string) => {
    localStorage.setItem(_props.name, showStr);

    setShow(+showStr);
    show() == 0 ? setShow(200) : '';
  };

  return (
    <label class="wy-inline-label">Bars per host
      <input type="number" min="1" class="wy-input" style={{ width: "7em" }} placeholder="200"
        value={localStorage.getItem(_props.name) || ""}
        onInput={e => handleSaveShow(e.currentTarget.value)}/>
    </label>
  )
}

export default HistShow
