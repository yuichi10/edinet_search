'use client';
import { FormEvent, useState, useRef } from "react";
import { GraphQLClient } from 'graphql-request';

type Company = {
  docID: string;
  secCode: string;
  filerName: string;
  docDescription: string;
  submitDatetime: string;
  avgAge: string;
  avgYearOfService: string;
  avgAnnualSalary: string;
  numberOfEployees: string;
  employeeInformation: string;
};

type Companies = {
    Companies: Company[];
};

export default function Home() {
  const [companies, setCompanies] = useState<Company[]>([]);

  const client = new GraphQLClient('http://localhost:8080/api/query')
  const query = `
    query SearchCompanies($filerName: String!) {
      Companies(filter: {filerName: $filerName}) {
        docID
        secCode
        filerName
        docDescription
        submitDatetime
        avgAge
        avgYearOfService
        avgAnnualSalary
        numberOfEmployees
        employeeInformation
      }
    }
  `;
  const submitSearch = async(event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const inputCompanyName = event.currentTarget.elements.namedItem("inputCompanyName") as HTMLInputElement;
    let variables = {
      filerName: ""
    }

    console.log(inputCompanyName.value)
    if (inputCompanyName.value != "") {
      variables.filerName = inputCompanyName.value;
    }

    try {
      const data: Companies = await client.request(query, variables);
      console.log(data)
      setCompanies(data.Companies)
    } catch(e) {
      console.log("failed to fetch company data")
      console.log(e)
    }
  }

  return (
    <div>
      <h1 className="font-sans font-bold p-px:4px">Edinet検索</h1>
      <div>
        <form className="max-w-sm mx-auto" onSubmit={submitSearch}>
          <div className="mb-5">
            <label htmlFor="text" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">会社名</label>
            <input type="text" name="inputCompanyName" id="filerName" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="株式会社" required />
          </div>
          <button type="submit" className="p-5 text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Submit</button>
        </form>
      </div>

      { companies.length > 0 &&
      <div className="p-10">
        <table className="table-auto border border-collapse border-slate-500">
          <thead>
            <tr>
              <th className="border border-slate-600">会社名</th>
              <th className="border border-slate-600">勤続年数</th>
              <th className="border border-slate-600">平均年齢</th>
              <th className="border border-slate-600">平均年収</th>
              <th className="border border-slate-600">従業員数</th>
              <th className="border border-slate-600">情報の追加日</th>
              {/* 他のヘッダー... */}
            </tr>
          </thead>
          <tbody>
            {companies.map((company) => (
              <tr key={company.docID}>
                <td className="border border-slate-700"><a href={`https://disclosure2dl.edinet-fsa.go.jp/searchdocument/pdf/${company.docID}.pdf`}>{company.filerName}</a></td>
                <td className="border border-slate-700">{company.avgYearOfService}年</td>
                <td className="border border-slate-700">{company.avgAge}歳</td>
                <td className="border border-slate-700">{company.avgAnnualSalary}円</td>
                <td className="border border-slate-700">{company.numberOfEployees}人</td>
                <td className="border border-slate-700">{company.submitDatetime}</td>
                {/* <td>{company.employeeInformation}</td> */}
                {/* 他のデータセル... */}
              </tr>
            ))}
          </tbody>
        </table>
        </div>
      }
    </div>
  )
}
