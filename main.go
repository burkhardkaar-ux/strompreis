package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type StromRechnung struct {
	VerbrauchKWh        float64
	EnergieArbeitspreis float64
	EnergieKostenNetto  float64
	EnergieGrundpreis   float64
	NetzArbeitspreis    float64
	NetzVerlustNetto    float64
	NetzGrundpreis      float64
	ElektroAbgabeNetto  float64
	NettoGesamt         float64
	Umsatzsteuer        float64
	BruttoGesamt        float64
	MonatsAbschlag      float64
	EffektiverKWhPreis  float64
}

type PageData struct {
	Berechnet   bool
	Fehler      string
	VerbrauchIn string
	PreisIn     string
	Ergebnis    StromRechnung
}

func berechneStrompreis(verbrauch float64, preisProKWh float64) StromRechnung {
	energieGrundpreis := 59.90
	netzArbeitspreisProKWh := 9.11
	netzVerlustProKWh := 0.40
	netzGrundpreis := 210.00
	elektroAbgabeProKWh := 1.50

	energieKostenNetto := verbrauch * (preisProKWh / 100.0)
	netzArbeitNetto := verbrauch * (netzArbeitspreisProKWh / 100.0)
	netzVerlustNetto := verbrauch * (netzVerlustProKWh / 100.0)
	abgabeNetto := verbrauch * (elektroAbgabeProKWh / 100.0)

	nettoGesamt := energieKostenNetto + energieGrundpreis + netzArbeitNetto + netzVerlustNetto + netzGrundpreis + abgabeNetto
	ust := nettoGesamt * 0.20
	bruttoGesamt := nettoGesamt + ust
	monatsAbschlag := bruttoGesamt / 12.0
	effektiverPreis := (bruttoGesamt / verbrauch) * 100.0

	return StromRechnung{
		VerbrauchKWh:        verbrauch,
		EnergieArbeitspreis: preisProKWh,
		EnergieKostenNetto:  energieKostenNetto,
		EnergieGrundpreis:   energieGrundpreis,
		NetzArbeitspreis:    netzArbeitNetto,
		NetzVerlustNetto:    netzVerlustNetto,
		NetzGrundpreis:      netzGrundpreis,
		ElektroAbgabeNetto:  abgabeNetto,
		NettoGesamt:         nettoGesamt,
		Umsatzsteuer:        ust,
		BruttoGesamt:        bruttoGesamt,
		MonatsAbschlag:      monatsAbschlag,
		EffektiverKWhPreis:  effektiverPreis,
	}
}

func main() {
	// Die beiden HTML-Teile sicher zusammenfügen
	vollstaendigesHTML := htmlPart1 + htmlPart2
	tmpl := template.Must(template.New("index").Parse(vollstaendigesHTML))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			VerbrauchIn: "",
			PreisIn:     "",
		}

		if r.Method == http.MethodPost {
			vStr := strings.TrimSpace(strings.ReplaceAll(r.FormValue("verbrauch"), ",", "."))
			pStr := strings.TrimSpace(strings.ReplaceAll(r.FormValue("preis"), ",", "."))

			data.VerbrauchIn = vStr
			data.PreisIn = pStr

			verbrauch, errV := strconv.ParseFloat(vStr, 64)
			preis, errP := strconv.ParseFloat(pStr, 64)

			if errV != nil || errP != nil || verbrauch <= 0 || preis <= 0 {
				data.Fehler = "Bitte gib gültige positive Zahlen für den Verbrauch und den Energiepreis ein."
			} else {
				data.Ergebnis = berechneStrompreis(verbrauch, preis)
				data.Berechnet = true
			}
		}

		tmpl.Execute(w, data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server läuft auf Port %s ...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

const htmlPart1 = `<!-- Copyright by Ing. Burkhard Kaar, 8580 Koeflach, Austria -->

<!-- Copyright by Ing. Burkhard Kaar, 8580 Koeflach, Austria -->

<!-- Copyright by Ing. Burkhard Kaar, 8580 Koeflach, Austria -->

<!-- Copyright by Ing. Burkhard Kaar, 8580 Koeflach, Austria -->

<!-- Copyright by Ing. Burkhard Kaar, 8580 Koeflach, Austria -->

<!-- Copyright by Ing. Burkhard Kaar, 8580 Koeflach, Austria -->

<!DOCTYPE html>
<html lang="de">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Berechnung Durchlaufträger</title>
    <style>
        body {
            font-family: 'Roboto', sans-serif;
            max-width: 800px;
            margin: 20px auto;
            padding: 20px;
            background: #000000;
            color: #333;
        }

        .card {
            background: white;
            padding: 25px;
            border-radius: 5px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
        }

        h1 {
            color: #111;
            border-bottom: 2px solid #eee;
            padding-bottom: 8px;
            margin-top: 0;
            font-size: 24px;
        }

        h2 {
            color: #111;
            margin-top: 20px;
            font-size: 20px;
            border-bottom: 1px solid #e2e8f0;
            padding-bottom: 4px;
        }

        h3 {
            color: #428787;
            margin-bottom: 12px;
            border-bottom: 1px solid #e2e8f0;
            padding-bottom: 4px;
            font-size: 18px;
        }

        .grid {
            display: grid;
            grid-template-columns: 1fr;
            gap: 20px;
            margin-bottom: 20px;
        }

        .box {
            background: #f8fafc;
            padding: 15px;
            border-radius: 5px;
            border: 1px solid #e2e8f0;
        }

        .form-group {
            margin-bottom: 15px;
        }

        label {
            display: block;
            font-size: 14px;
            font-weight: bold;
            color: #475569;
            margin-bottom: 6px;
        }

        input[type="number"],
        input[type="text"] {
            width: 100%;
            padding: 8px;
            border: 1px solid #cbd5e1;
            font-size: 16px;
            box-sizing: border-box;
        }

        button {
            background: #428787;
            color: white;
            border: none;
            padding: 12px;
            cursor: pointer;
            font-size: 16px;
            border-radius: 4px;
            width: 100%;
            font-weight: bold;
            transition: background 0.2s;
        }

        button:hover {
            background: #336666;
        }

        .error {
            color: #b91c1c;
            font-weight: bold;
            background: #fee2e2;
            padding: 12px;
            border-radius: 4px;
            border: 1px solid #fca5a5;
            margin-bottom: 20px;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 15px;
            margin-bottom: 20px;
            font-size: 14px;
        }

        th,
        td {
            border: 1px solid #e2e8f0;
            padding: 8px 12px;
            text-align: left;
        }

        th {
            background-color: #428787 !important;
            color: white !important;
        }

        tr:nth-child(even) {
            background-color: #f8fafc;
        }

        .r {
            text-align: right;
        }

        .hint {
            color: #64748b;
            font-size: 12px;
            margin-top: 4px;
            display: block;
            font-weight: normal;
        }

    </style>
</head>


<body>
    <div class="card">
        <h1>Berechnung Durchlaufträger</h1>
        <p style="color: #428787; font-style: italic; margin-top: -5px; margin-bottom: 20px; font-size: 14px;">Systembelastung: Konstante Gleichlast q = 10.0 kN/m</p>

        {{if .Fehler}}
        <div class="error">{{.Fehler}}</div>
        {{end}}

        <form method="POST" id="traegerForm">
            <div class="grid">
                <div class="box">
                    <h3>Eingabeparameter Träger</h3>
                    <div class="form-group">
                        <label for="felder">Anzahl der Felder:</label>
                        <input type="number" id="felder" name="felder" placeholder="3" value="{{if .Berechnet}}{{.AnzahlFelder}}{{end}}" min="1" required>
                    </div>
                    <div class="form-group">
                        <label for="laengen">Feldlängen in Metern (mit Komma getrennt):</label>
                        <input type="text" id="laengen" name="laengen" value="{{.LaengenInput}}" placeholder="z.B. 4.5, 6.0, 4.5" required>
                        <span class="hint">Gib die Längen der Felder der Reihe nach ein (z.B. 5, 4.5, 5)</span>
                    </div>

                    <div>
                        <h6>&copy; Ing. Burkhard Kaar - Studienprojekt imperative Programmierung, Berechnung eines Durchlaufträgers nach dem Cross-Verfahren</h6>
                    </div>

                    <button type="submit" id="submitBtn" disabled>Berechnung starten</button>

                    {{if .Berechnet}}
                    <h2 class="result-title">Berechnungsergebnisse</h2>
                    <table>
                        <thead>
                            <tr>
                                <th>Stabseite</th>
                                <th class="r">µ</th>
                                <th class="r">Moment (kNm)</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .Seiten}}
                            <tr>
                                <td><strong>{{.Name}}</strong></td>
                                <td class="r">{{printf "%.2f" .Mu}}</td>
                                <td class="r">{{printf "%.2f" .MTotal}}</td>
                            </tr>
                            {{end}}
                        </tbody>
                    </table>
                    {{end}}
                </div>
            </div>
        </form>

        <style>
            /* Standard-Design für den aktiven Button */
            button {
                background: #428787;
                color: white;
                border: none;
                padding: 12px;
                cursor: pointer;
                font-size: 16px;
                border-radius: 4px;
                width: 100%;
                font-weight: bold;
                transition: all 0.2s;
            }

            button:hover:not(:disabled) {
                background: #336666;
            }

            /* Automatisches Ausgrauen, wenn der Button deaktiviert ist */
            button:disabled {
                background: #cbd5e1;
                /* Ein helles Grau */
                color: #94a3b8;
                /* Dunklerer Text für Lesbarkeit im Grau */
                cursor: not-allowed;
                /* Zeigt ein "Verboten"-Symbol beim Drüberfahren */
                opacity: 0.7;
            }

        </style>

        <!-- Dieser kleine JavaScript-Block steuert die Freischaltung live beim Tippen -->
        <script>
            const form = document.getElementById('traegerForm');
            const submitBtn = document.getElementById('submitBtn');

            function pruefeFormular() {
                // checkValidity() prüft automatisch alle Felder mit "required"
                if (form.checkValidity()) {
                    submitBtn.removeAttribute('disabled');
                } else {
                    submitBtn.setAttribute('disabled', 'true');
                }
            }

            // Lauscht auf jede Eingabe in den Formularfeldern
            form.addEventListener('input', pruefeFormular);

            // Einmal beim Laden ausführen (falls nach Berechnung bereits Werte drinstehen)
            pruefeFormular();

            function pruefeFormular() {
                var form = document.getElementById('stromForm');
                var submitBtn = document.getElementById('submitBtn');
                if (form && submitBtn) {
                    if (form.checkValidity()) {
                        submitBtn.removeAttribute('disabled');
                    } else {
                        submitBtn.setAttribute('disabled', 'true');
                    }
                }
            }

            document.addEventListener('DOMContentLoaded', function() {
                var form = document.getElementById('stromForm');
                if (form) {
                    form.addEventListener('input', pruefeFormular);
                }
                pruefeFormular();

                // DIESER TEIL SCHIEBT DAS FORMULAR NACH OBEN, WENN DIE ERGEBNISSE DA SIND:
                // Der Server fügt bei erfolgreicher Berechnung das <h2>-Element ein.
                // Sobald dieses Element im Browser existiert, scrollt die Seite automatisch dorthin.
                var ergebnisUeberschrift = document.querySelector('h2');
                if (ergebnisUeberschrift) {
                    ergebnisUeberschrift.scrollIntoView({
                        behavior: 'smooth',
                        block: 'start'
                    });
                }
            });

        </script>

    </div>
</body>

</html>

`
